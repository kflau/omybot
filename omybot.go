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
	token := ""
	flag.StringVar(&token, "token", "", "Discord Bot token")
	flag.StringVar(&webhook, "webhook", "", "Discord Webhook URL")
	flag.StringVar(&hkexToken, "hkexToken", "", "HKEX Token")
	flag.Parse()
	if token == "" || webhook == "" || hkexToken == "" {
		flag.PrintDefaults()
		return
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Printf("Error creating Discord session: %v\n", err)
		return
	}

	dg.AddHandler(ready)
	dg.AddHandler(connect)
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
	args := strings.Fields(m.Content)[1:]
	if len(args) <= 0 {
		fmt.Printf("Please specify RIC\n")
		return
	}
	ric := strings.TrimSpace(args[0])
	quote := &Quote{webhook: webhook, hkexToken: hkexToken, ric: ric}
	if err := quote.getQuote(); err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	if err := quote.sendQuote(); err != nil {
		fmt.Printf("%v\n", err)
		return
	}
}

func connect(s *discordgo.Session, m *discordgo.Connect) {

}