package main

import (
	"log"
	"time"

	"github.com/jvikstedt/bluestorm/hub"
	"github.com/jvikstedt/bluestorm/network"
)

type Lobby struct {
	*hub.BaseRoom
}

func NewLobby() *Lobby {
	lobby := &Lobby{}
	lobby.BaseRoom = hub.NewBaseRoom("lobby")
	return lobby
}

func (l *Lobby) Run() {
	for {
		time.Sleep(time.Second * 1)

		l.GetUsersWithRead(func(users hub.Users) {
			log.Printf("Room lobby users: %v\n", users)
		})
	}
}

func (l *Lobby) NewMsg(agent *network.Agent, msg interface{}) {
	l.Broadcast(msg)
}
