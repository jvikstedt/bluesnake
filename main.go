package main

import (
	"log"
	"os"

	"github.com/jvikstedt/bluestorm"
	"github.com/jvikstedt/bluestorm/hub"
	"github.com/jvikstedt/bluestorm/network"
	"github.com/jvikstedt/bluestorm/network/json"
)

type Msg struct {
	Message string `json:"name"`
}

func proxy(agent *network.Agent, i interface{}) {
	r, err := agent.GetValue("room")
	if err != nil {
		log.Println(err)
		return
	}

	room, ok := r.(hub.Room)
	if !ok {
		log.Println("Could not convert to room...")
		return
	}

	room.NewMsg(agent, i)
}

func main() {
	processor := json.NewProcessor()
	processor.Register(&Msg{}, proxy)

	lobby := NewLobby()
	hub.DefaultManager().AddRoom(lobby)

	servers := bluestorm.Servers{
		&network.WSServer{
			Addr:         ":8081",
			Processor:    processor,
			OnConnect:    bluestorm.OnConnectHelper(hub.DefaultManager(), lobby.ID()),
			OnDisconnect: bluestorm.OnDisconnectHelper(hub.DefaultManager()),
		},
	}

	go lobby.Run()

	bluestorm.Run(bluestorm.CloseOnSignal(os.Interrupt), servers)
}
