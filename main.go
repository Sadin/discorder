package main

import (
	"fmt"
	"os"
	"flag"
    	"database/sql"
    	_ "github.com/godror/godror"
    	"github.com/godror/godror"
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
	/// register guildUpdate event
	dg.AddHandler(guildUpdate)
	// register message create event handler
	dg.AddHandler(messageCreate)
	// register message edit event handler
	dg.AddHandler(messageUpdate)

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
	// check if guild is already seen in database
	if checkGuild(event.Guild.ID, event.Guild.Name) == 0 {
	    // store guild
	    go storeGuild(event.Guild.ID, event.Guild.Name)
	}

	guild_list = append(guild_list,event.Guild.Name)
}

func guildUpdate(s *discordgo.Session, event *discordgo.GuildUpdate) {
	fmt.Println(event.Guild.Name, "updated")
}

func checkGuild(id string, name string) int {
	var count int
	sql := fmt.Sprintf("SELECT count(*) FROM guilds WHERE guild_id = %s", id)
	row := DB.QueryRow(sql, godror.FetchArraySize(1))

	row.Scan(&count)

	if count == 0 {
	    fmt.Println(name, "not found in db...")
	} else {
	    fmt.Println(name, "found in db...")
	}

	return count
}

func storeGuild(id string, name string) {
	_, err := DB.Exec("INSERT INTO guilds VALUES (:1, :2)", id, name)
	if err != nil {
	    fmt.Println(".....Error Inserting guild data")
	    fmt.Println(err)
	    return
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore messages sent by bot
	if m.Author.ID == s.State.User.ID {
		return
	}
	// insert in goroutine
	go logMessage(s, m)

	message := fmt.Sprintf(m.Content)
	fmt.Println(message)
}

func messageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	//update information in message_log
	_, err := DB.Exec("UPDATE message SET message_content = :1 WHERE message_id = :2 and message_guild_id = :3 and message_channel_id = :4", m.Content, m.ID, m.GuildID, m.ChannelID)
	if err != nil {
	    fmt.Println(".....Error updating message data")
	    fmt.Println(err)
	    return
	}
	fmt.Println(fmt.Sprintf(`Message %s updated in %s`, m.ID, m.GuildID))
}

func logMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := DB.Exec("INSERT INTO message VALUES (:1, :2, :3, :4, :5, :6, :7)", m.ID, m.Timestamp, m.GuildID, m.ChannelID, m.Author.ID, fmt.Sprintf("%s",m.Author), m.Content)
	if err != nil {
	    fmt.Println(".....Error Inserting message data")
	    fmt.Println(err)
	    return
	}
}

