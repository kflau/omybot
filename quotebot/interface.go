package main

import (
    "github.com/bwmarrin/discordgo"
)

type MessageCreateHandler interface {

	New() *MessageCreateHandler

	Type() string

    Parse([]string) bool

    String() string

	MemberJoin(*discordgo.MessageCreate) (string, error)

	Forward([]string) (error)

	Reply([]string) (error)
}
