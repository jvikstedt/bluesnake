package main

import (
	"log"
	"time"

	"github.com/jvikstedt/bluestorm/hub"
)

type Lobby struct {
	*hub.BaseRoom

	msgs chan *Msg
}

func NewLobby(id hub.RoomID) *Lobby {
	return &Lobby{
		BaseRoom: hub.NewBaseRoom(id),
		msgs:     make(chan *Msg, 10),
	}
}

func (l *Lobby) Run() {
	for {
		time.Sleep(time.Second * 1)

		select {
		case msg := <-l.msgs:
			msg.Run(l, msg)
		default:
			break
		}

		l.GetUsersWithRead(func(users hub.Users) {
			log.Printf("Room lobby users: %v\n", users)
		})
	}
}

func (l *Lobby) AddMsg(msg *Msg) {
	l.msgs <- msg
}
