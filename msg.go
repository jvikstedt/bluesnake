package main

type Command interface {
	Run(receiver Receiver, msg *Msg)
}

type Receiver interface {
	AddMsg(msg *Msg)
}

type Msg struct {
	*Player
	Command
}
