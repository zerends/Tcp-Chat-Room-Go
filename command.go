package main

type commandID int

const (
	cmdNick commandID = iota
	cmdJoin
	cmdList
	cmdRoom
	cmdMsg
	cmdQuit
	cmdHelp
)

type command struct {
	id     commandID
	client *client
	args   []string
}
