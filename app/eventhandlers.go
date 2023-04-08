package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var GuildVoiceConnections = make(map[string]string)

func RegisterEvents() {
	// Ready event.
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		intents := s.Identify.Intents

		fmt.Printf("Logged in as: %v#%v\nIntents: %v\n", s.State.User.Username, s.State.User.Discriminator, intents)
	})

	// Command handler.
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := CommandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
			return
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Invalid command! Deleting...",
			},
		})

		commandID := i.Interaction.ApplicationCommandData().ID

		fmt.Printf("Invalid command detected.\nCommand ID: %s\nCommand used by: %s", commandID, i.Interaction.Member.User.ID)
		s.ApplicationCommandDelete(s.State.User.ID, "", commandID)
	})

	// Presence update event.
	s.AddHandler(func(s *discordgo.Session, p *discordgo.PresenceUpdate) {
		guildID, ok := GuildVoiceConnections[p.User.ID]
		if !ok {
			return
		}

		memberState, err := s.State.VoiceState(guildID, p.User.ID)
		if err != nil {
			fmt.Println(err)
			return
		}

		botPermissions, err := s.UserChannelPermissions(s.State.User.ID, memberState.ChannelID)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Check if the bot has permissions to move member.
		if botPermissions&discordgo.PermissionVoiceMoveMembers == 0 {
			return
		}

		var blacklistedActivity string
		guildData := getGuildData(p.GuildID)

		for _, activity := range p.Activities {
			if len(activity.ApplicationID) <= 0 {
				continue
			}

			if isInArray(activity.ApplicationID, guildData.Blacklists) {
				activityName, ok := Activities[activity.ApplicationID]

				if !ok {
					continue
				}

				blacklistedActivity = activityName
				break
			}
		}

		if len(blacklistedActivity) > 0 {
			s.GuildMemberMove(guildID, p.User.ID, nil)
			private, err := s.UserChannelCreate(p.User.ID)
			if err != nil {
				return
			}

			message := fmt.Sprintf(DisconnectMessage, p.User.ID, memberState.ChannelID, blacklistedActivity)
			_, err = s.ChannelMessageSend(private.ID, message)
			if err != nil {
				fmt.Println(err)
			}
		}
	})

	// Voice state change.
	s.AddHandler(func(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
		if len(v.ChannelID) <= 0 {
			delete(GuildVoiceConnections, v.UserID)
			return
		}

		GuildVoiceConnections[v.UserID] = v.GuildID
	})
}
