package main

type MessageCreateHandler interface {

	New() interface{}

	Type() string

	MemberJoin(*discordgo.MessageCreate) (string, error)

	Forward(*discordgo.MessageCreate) (string, error)

	Reply(*discordgo.MessageCreate) (string, error)
}