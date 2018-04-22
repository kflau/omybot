package main

import (
	"fmt"
	"strings"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"github.com/bwmarrin/discordgo"
)

var (
	webhook 		string
	botName 		string
)

func main() {
	token := ""
	flag.StringVar(&token, "token", "", "Discord Bot token")
	flag.StringVar(&webhook, "webhook", "", "Discord Webhook URL")
	flag.StringVar(&botName, "botName", "", "Bot Name")
	flag.Parse()
	if token == "" || webhook == "" || botName == "" {
		flag.PrintDefaults()
		return
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Printf("Error creating Discord session: %v\n", err)
		return
	}

	dg.AddHandler(ready)
	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Printf("Error opening Discord session: %v\n", err)
	}

	fmt.Printf("QuoteBot is now running.  Press CTRL-C to exit.\n")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateStatus(0, "QuoteBot")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Messages are boardcasted even to myself, ignore message sent by myself
	if botName == m.Author.Username && m.Author.Bot {
		return
	}

	qt := NewQuote()
	if m.Type == discordgo.MessageTypeGuildMemberJoin {
		if str, err := qt.MemberJoin(m); err != nil {
			fmt.Printf("Cannot handle MemberJoin [%#v] due to %#v\n", str, err)
		}
		return
	}
	if len(m.Content) <= 0 {
		return
	}
	args := strings.Fields(m.Content)[1:]
	if !qt.Parse(args) {
		fmt.Printf("Not a quote format\n")
		return
	}
	if err := qt.Forward(args); err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	args = append([]string{webhook}, args...)
	if err := qt.Reply(args); err != nil {
		fmt.Printf("%v\n", err)
		return
	}
}
