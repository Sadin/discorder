package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"github.com/bwmarrin/discord.go"
)

// system variables
var (
	Token string
	buffer = make([][]byte, 0)
)
// global variables
var guild_list []string

func main() {
	fmt.Println("discorder")

	// launch discord session
	fmt.Print("creating new discord session...")
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register guildCreate as a callback for the guildCreate events.
	dg.AddHandler(guildCreate)
	// register message create event handler
	dg.AddHandler(messageCreate)
}

func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {

	if event.Guild.Unavailable {
		fmt.Println("error with guild" + event.Guild.ID)
		return
	}
	guild_list = append(guild_list,event.Guild.Name)
	fmt.Println(event.Guild.Name)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore messages sent by bot
	if m.Author.ID == s.State.User.ID {
		return
	}

	fmt.Println("Message detected")
}