package main

import (
	"fmt"
	"os"
	"flag"
    	"database/sql"
    	_ "github.com/godror/godror"
	"os/signal"
	"syscall"
	"github.com/bwmarrin/discordgo"
)

// system variables
var (
	Token string
	Dbuser string
	Dbpass string
	buffer = make([][]byte, 0)

	DB *sql.DB
)
// global variables
var guild_list []string

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&Dbuser, "u", "", "DB username")
	flag.StringVar(&Dbpass, "p", "", "DB Pass" )
	flag.Parse()
}

func main() {
	fmt.Println("discorder")

	// database
	connectDB()

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

func connectDB () {
	fmt.Println("connecting to database...")
	db, err := sql.Open("godror", fmt.Sprintf(`user="%s" password="%s" connectString="localhost:1521"`, Dbuser, Dbpass))
    	if err != nil {
        	fmt.Println(err)
        	return
    	}

    	DB = db
}

func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {

	if event.Guild.Unavailable {
		fmt.Println("error with guild" + event.Guild.ID)
		return
	}
	_, err := DB.Exec("INSERT INTO guilds VALUES (:1, :2)",event.Guild.ID, event.Guild.Name)
	if err != nil {
	    fmt.Println(".....Error Inserting guild data")
	    fmt.Println(err)
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

	go logMessage(s, m)

	message := fmt.Sprintf(m.Content)
	fmt.Println(message)
}

func logMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := DB.Exec("INSERT INTO message VALUES (:1, :2, :3, :4, :5)", m.ID, m.Timestamp, m.Author.ID, fmt.Sprintf("%s",m.Author), m.Content)
	if err != nil {
	    fmt.Println(".....Error Inserting message data")
	    fmt.Println(err)
	    return
	}
}
