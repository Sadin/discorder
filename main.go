package main

import (
	"fmt"
	"os"
	"flag"
	"os/signal"
	"syscall"
	"github.com/bwmarrin/discordgo"
)

// system variables
var (
	Token string
	buffer = make([][]byte, 0)
)
// global variables
var guild_list []string

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	fmt.Println("discorder")

	// launch discord session
	fmt.Println("creating new discord session...")
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register guildCreate as a callback for the guildCreate events.
	dg.AddHandler(guildCreate)
	// register message create event handler
	dg.AddHandler(messageCreate)

	// open connection to discord
	fmt.Print("connecting to discord")
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// wait for CTRL-C or term sig
	fmt.Println("Punx operational. Press CTRL-C to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// close
	dg.Close()
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
