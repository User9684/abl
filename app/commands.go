package main

import (
	"github.com/bwmarrin/discordgo"
)

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "blacklist",
		Description: "Blacklist a specific activity",
		Options: []*discordgo.ApplicationCommandOption{

			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "activity",
				Description: "Activity to blacklist",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Putt Party",
						Value: "945737671223947305",
					},
				},
			},
		},
	},
	{
		Name:        "getvoicestate",
		Description: "Get the voice state of a specific member in a guild",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "guild",
				Description: "Guild ID",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user",
				Description: "User ID",
				Required:    true,
			},
		},
	},
	{
		Name:        "getuserpresence",
		Description: "Get the target users current presence, if found.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "guild",
				Description: "Guild ID",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user",
				Description: "User ID",
				Required:    true,
			},
		},
	},
}
