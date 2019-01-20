package main

type Command interface {
	Run(receiver Receiver, player *Player)
}

type Receiver interface {
	AddMsg(msg *Msg)
}

type Msg struct {
	player *Player
	cmd    Command
}
