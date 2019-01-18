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
	msg, ok := i.(*Msg)
	if !ok {
		log.Println("not right type")
		return
	}

	r, err := agent.GetValue("room")
	if err != nil {
		log.Println(err)
		return
	}

	room, ok := r.(*hub.Room)
	if !ok {
		log.Println("not okay...")
		return
	}

	room.BroadcastExceptOne(hub.UserID(agent.ID()), msg)
}

func main() {
	processor := json.NewProcessor()
	processor.Register(&Msg{}, proxy)

	servers := bluestorm.Servers{
		&network.WSServer{
			Addr:         ":8081",
			Processor:    processor,
			OnConnect:    bluestorm.OnConnectHelper(hub.DefaultManager(), hub.DefaultRoomID),
			OnDisconnect: bluestorm.OnDisconnectHelper(hub.DefaultManager()),
		},
	}

	bluestorm.Run(bluestorm.CloseOnSignal(os.Interrupt), servers)
}
