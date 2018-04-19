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
	webhook string 		= ""
	hkexToken string 	= ""
)

func main() {
	tokenPtr := flag.String("token", "", "Discord Bot token")
	webhookPtr := flag.String("webhook", "", "Discord Webhook URL")
	hkexTokenPtr := flag.String("hkexToken", "", "HKEX Token")
	flag.Parse()

	if *tokenPtr == "" || *webhookPtr == "" || *hkexTokenPtr == "" {
		flag.PrintDefaults()
		return
	}

	token := *tokenPtr
	webhook = *webhookPtr
	hkexToken = *hkexTokenPtr
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

	fmt.Printf("OMyBot is now running.  Press CTRL-C to exit.\n")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateStatus(0, "OMyBot")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Messages are boardcasted even to myself, ignore message sent by myself
	if m.Author.ID == s.State.User.ID {
		return
	}
	fmt.Printf("Author: %v, User: %v, Received Message: %v\n", m.Author.ID, s.State.User.ID, m.Content)
	if strings.HasPrefix(m.Content, "!quote ") || strings.HasPrefix(m.Content, "!q ") {
		args := strings.Fields(m.Content)[1:]
		quote := new(Quote)
		quote.webhook = webhook
		quote.hkexToken = hkexToken
		if err := quote.getQuote(args); err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		if err := quote.sendQuote(); err != nil {
			fmt.Printf("%v\n", err)
			return
		}
	}
}
