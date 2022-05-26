package main

import (
	"fmt"
	"os"
	"strings"
	"regexp"
	"flag"
    	"database/sql"
    	_ "github.com/godror/godror"
    	"github.com/godror/godror"
	"os/signal"
	"syscall"
	"github.com/bwmarrin/discordgo"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// system variables
var (
	Token string
	Dbuser string
	Dbpass string
	buffer = make([][]byte, 0)

	DB *sql.DB

	// logging
	logger log.Logger
)

// global variables
var guild_list []string
var registered_commands = make(map[string]int)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&Dbuser, "u", "", "DB username")
	flag.StringVar(&Dbpass, "p", "", "DB Pass" )
	flag.Parse()
}

func main() {
	fmt.Println("discorder")

	// populate commands map
	registered_commands["help"] = 1

	// set up logging
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = level.NewFilter(logger, level.AllowInfo())
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	// database
	connectDB()

	// launch discord session
	level.Info(logger).Log("msg", "creating new discord session")
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		logger.Log("err", fmt.Sprintf("error creating Discord session,", err))
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
	level.Info(logger).Log("msg", "connecting to discord")
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// wait for CTRL-C or term sig
	level.Info(logger).Log("msg", "bot online")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// close
	dg.Close()
}

func connectDB () {
	level.Info(logger).Log("msg", "connecting to database")
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
	level.Info(logger).Log("msg", fmt.Sprintf(event.Guild.Name, "updated"))
}

func checkGuild(id string, name string) int {
	var count int
	sql := fmt.Sprintf("SELECT count(*) FROM guilds WHERE guild_id = %s", id)
	row := DB.QueryRow(sql, godror.FetchArraySize(1))

	row.Scan(&count)

	if count == 0 {
	    level.Info(logger).Log("msg", fmt.Sprintf("%s not found in db...", name))
	} else {
	    level.Info(logger).Log("msg", fmt.Sprintf("%s found in db...", name))
	}

	return count
}

func storeGuild(id string, name string) {
	_, err := DB.Exec("INSERT INTO guilds VALUES (:1, :2)", id, name)
	if err != nil {
	    level.Warn(logger).Log("warn", ".....Error Inserting guild data")
	    level.Error(logger).Log("sql", err)
	    return
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore messages sent by bot
	if (m.Author.ID == s.State.User.ID) || (m.Author.Bot == true){
		return
	}
	// insert in goroutine
	go logMessage(s, m)

	// if command, parse and execute
	cmdStr, err := regexp.MatchString(`>>`, m.Content)
	if err != nil {
		level.Error(logger).Log("err", err)
	}
	if cmdStr == true {
		feedback := parseCommand(m)
		_,err := s.ChannelMessageSend(m.ChannelID, feedback)
		if err != nil {
			level.Error(logger).Log("err", err)
		}
	} else {
		level.Info(logger).Log("chat", fmt.Sprintf(m.Content))
		for i:=0; i<len(m.Attachments); i++ {
			if i==0 {
				level.Info(logger).Log("msg", "Attachment(s) Found...")
			}
			go parseAttachement(m.ID, m.Attachments[i])
		}
	}
}

func messageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	//update information in message_log
	_, err := DB.Exec("UPDATE message SET message_content = :1 WHERE message_id = :2 and message_guild_id = :3 and message_channel_id = :4", m.Content, m.ID, m.GuildID, m.ChannelID)
	if err != nil {
	    level.Warn(logger).Log(".....Error updating message data")
	    level.Error(logger).Log("sql", err)
	    return
	}
	level.Info(logger).Log("msg", fmt.Sprintf(`Message %s updated in %s`, m.ID, m.GuildID))
}

func parseAttachement (mid string, a *discordgo.MessageAttachment) {
	level.Info(logger).Log("attachment", fmt.Sprintf("ID: %s\n URL: %s\n ProxyURL: %s\n Filename %s\n ContentType %s\n Dimensions: %dx%d\n Size: %d", a.ID, a.URL, a.ProxyURL, a.Filename, a.ContentType, a.Width, a.Height, a.Size))
}

func parseCommand(m *discordgo.MessageCreate) string {
	// regex to clear command prefix
	re := regexp.MustCompile(`>>(.*)`)
	// split command from params substring, apply regex fitler to command
	message := strings.SplitN(m.Content, " ", 2)
	command, args := re.FindStringSubmatch(message[0]), message[1]

	level.Info(logger).Log("sys", "command received", "cmd", fmt.Sprintf(command[1]))
	if validateCommand(command[1]) {
		var empty []string
		level.Info(logger).Log("sys", "command exists, executing.")
		if message[1] != "" {
			return executeCommand(command[1], strings.Split(args, " "))
		} else {
			return executeCommand(command[1], empty)
		}
	} else {
		level.Info(logger).Log("sys", "command does not exist, informing user.")
		return "command not found."
	}
}

func validateCommand(cmd string) bool {
	// check if command is in the registered list
	fmt.Println(cmd)
	if _, ok := registered_commands[cmd]; ok {
		return true
	}
	// if "a" does not exist
	return false
}

func executeCommand(cmd string, params []string) string {
	fmt.Println(cmd)
	// run specified command function
	switch cmd {
	case "help":
		return executeHelpCommand()
	default:
		return "testing"
	}
}

func executeHelpCommand() string {
	return "help placeholder"
}

func logMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := DB.Exec("INSERT INTO message VALUES (:1, :2, :3, :4, :5, :6, :7)", m.ID, m.Timestamp, m.GuildID, m.ChannelID, m.Author.ID, fmt.Sprintf("%s",m.Author), m.Content)
	if err != nil {
	    level.Warn(logger).Log("warn", ".....Error Inserting message data")
	    level.Error(logger).Log("sql", err)
	    return
	}
}


