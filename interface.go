package main

type DiscordLifecycle interface {

	New() interface{}

	Type() string

	MemberJoin([]interface{}) (string, error)

	Forward([]interface{}) (string, error)

	Reply([]interface{}) (string, error)
}