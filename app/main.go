package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
)

var DisconnectMessage = "<@!%s> you have been disconnected from <#%s> because the Activity `%s` has been blacklisted by the server staff!"

var s *discordgo.Session
var c *mongo.Client
var d *mongo.Database
var Activities map[string]string

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

	// Register activity map.
	file, err := os.Open("/activities.json")
	if err != nil {
		fmt.Printf("Failed to open activities file! %s\n", err)
	}

	json.NewDecoder(file).Decode(&Activities)

	file.Close()

	// Bot setup.
	fmt.Println("Starting the bot...")
	BotInit()

	err = s.Open()
	if err != nil {
		fmt.Printf("Cannot open the session: %v\n", err)
		return
	}

	CmdInit(s)

	// Waits for SIGTERM.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}

func isInArray(value string, array []string) bool {
	for _, v := range array {
		if value == v {
			return true
		}
	}

	return false
}
