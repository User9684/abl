package main

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
)

var DisconnectMessage = "<@!%s> you have been disconnected from <#%s> because the Activity `%s` has been blacklisted by the server staff!"

var s *discordgo.Session
var c *mongo.Client
var d *mongo.Database

func BotInit() {
	// Create bot client.
	session, err := discordgo.New(os.Getenv("TOKEN"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	// Register intents.
	session.Identify.Intents |= discordgo.IntentGuilds
	session.Identify.Intents |= discordgo.IntentGuildPresences
	session.Identify.Intents |= discordgo.IntentGuildMembers
	session.Identify.Intents |= discordgo.IntentGuildVoiceStates

	// Setup state.
	session.StateEnabled = true
	session.State.TrackChannels = true
	session.State.TrackMembers = true
	session.State.TrackVoice = true
	session.State.TrackPresences = true

	s = session

	RegisterEvents()
}

func main() {
	// Mongo setup.
	fmt.Println("Connecting to mongo...")
	c = DbInit()
	d = c.Database("ABL")
	fmt.Println("Verifying all collections...")
	CollectionCheck(d)

	// Bot setup.
	fmt.Println("Starting the bot...")
	BotInit()

	err := s.Open()
	if err != nil {
		fmt.Printf("Cannot open the session: %v\n", err)
		return
	}

	CmdInit(s)

	// Prevent the process from terminating.
	for {

	}
}

func isInArray(value string, array []string) bool {
	for _, v := range array {
		if value == v {
			return true
		}
	}

	return false
}
