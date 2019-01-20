package main

import (
	"log"
	"time"

	"github.com/jvikstedt/bluestorm/hub"
)

type Lobby struct {
	*hub.BaseRoom
}

func NewLobby(id hub.RoomID) *Lobby {
	return &Lobby{
		BaseRoom: hub.NewBaseRoom(id),
	}
}

func (l *Lobby) Run() {
	for {
		time.Sleep(time.Second * 1)

		l.GetUsersWithRead(func(users hub.Users) {
			log.Printf("Room lobby users: %v\n", users)
		})
	}
}

func (l *Lobby) NewMsg(player *Player, msg interface{}) {
	l.BroadcastExceptOne(player.ID(), msg)
}
