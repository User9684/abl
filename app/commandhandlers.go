package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var ActivityCommandDefaultChoices = make([]*discordgo.ApplicationCommandOptionChoice, 0, 25)

var CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"rep": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// Convert options to map
		optionMap := mapOptions(i.ApplicationCommandData().Options)

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("%v", optionMap["text"].Value),
			},
		})
	},
	"blacklist": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			// Convert options to map
			optionMap := mapOptions(i.ApplicationCommandData().Options)

			if err := blacklistActivity(i.GuildID, fmt.Sprintf("%v", optionMap["activity"].Value)); err != nil {
				fmt.Println(err)
				cmdError(i, err)
				return
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("You selected:  `%v`", optionMap["activity"].Value),
				},
			})
		case discordgo.InteractionApplicationCommandAutocomplete:
			activityCommandAutocomplete(i)
		}
	},
	"getvoicestate": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// Convert options to map
		optionMap := mapOptions(i.ApplicationCommandData().Options)
		guildID := fmt.Sprintf("%v", optionMap["guild"].Value)
		userID := fmt.Sprintf("%v", optionMap["user"].Value)

		memberState, err := s.State.VoiceState(guildID, userID)
		if err != nil {
			cmdError(i, err)
			return
		}

		j, err := json.MarshalIndent(memberState, "", "	")
		if err != nil {
			cmdError(i, err)
			return
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: string(j),
			},
		})
	},
	"getuserpresence": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// Convert options to map
		optionMap := mapOptions(i.ApplicationCommandData().Options)
		guildID := fmt.Sprintf("%v", optionMap["guild"].Value)
		userID := fmt.Sprintf("%v", optionMap["user"].Value)

		presence, err := s.State.Presence(guildID, userID)
		if err != nil {
			cmdError(i, err)
			return
		}

		j, err := json.MarshalIndent(presence.Activities, "", "	")
		if err != nil {
			cmdError(i, err)
			return
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: string(j),
			},
		})
	},
}

func CmdInit(s *discordgo.Session) {
	registeredCommands := make([]*discordgo.ApplicationCommand, len(Commands))
	for i, v := range Commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			fmt.Printf("Cannot create '%v' command: %v\n", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	for id, name := range Activities {
		choice := &discordgo.ApplicationCommandOptionChoice{
			Name:  name,
			Value: id,
		}

		ActivityCommandDefaultChoices = append(ActivityCommandDefaultChoices, choice)

		if len(ActivityCommandDefaultChoices) >= 25 {
			break
		}
	}
}

// Autocomplete function for activity based commands.
func activityCommandAutocomplete(i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	choices := make([]*discordgo.ApplicationCommandOptionChoice, 0, 25)
	input := strings.ToLower(data.Options[0].StringValue())

	for id, name := range Activities {
		if len(choices) >= 25 {
			break
		}

		choice := &discordgo.ApplicationCommandOptionChoice{
			Name:  name,
			Value: id,
		}

		if strings.Contains(strings.ToLower(name), input) {
			choices = append(choices, choice)
		}
	}

	if len(input) <= 1 {
		choices = ActivityCommandDefaultChoices
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})

	if err != nil {
		fmt.Println(err)
	}
}

func mapOptions(options []*discordgo.ApplicationCommandInteractionDataOption) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	var optionMap = make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	return optionMap
}

func cmdError(i *discordgo.InteractionCreate, err error) {
	if err == nil {
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("An error occured! \n```%s```", err.Error()),
		},
	})
}
