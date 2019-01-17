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

func onConnect(agent *network.Agent) {
	log.Printf("Agent connectedted %s\n", agent.ID())
	defaultRoom, _ := hub.DefaultManager().GetRoom(hub.DefaultRoomID)
	err := hub.DefaultManager().UserToRoom(hub.UserID(agent.ID()), hub.DefaultRoomID, agent)
	if err != nil {
		log.Println(err)
		agent.Conn().Close()
		return
	}
	agent.SetValue("room", defaultRoom)
	log.Printf("room has users: %v\n", defaultRoom.GetUsers())
}

func onDisconnect(agent *network.Agent) {
	log.Printf("Agent disconnected %s\n", agent.ID())
	defaultRoom, _ := hub.DefaultManager().GetRoom(hub.DefaultRoomID)
	err := hub.DefaultManager().RemoveUser(hub.UserID(agent.ID()))
	if err != nil {
		log.Println(err)
	}
	log.Printf("room has users: %v\n", defaultRoom.GetUsers())
}

func main() {
	processor := json.NewProcessor()
	processor.Register(&Msg{}, proxy)

	servers := bluestorm.Servers{
		&network.WSServer{
			Addr:         ":8081",
			Processor:    processor,
			OnConnect:    onConnect,
			OnDisconnect: onDisconnect,
		},
	}

	bluestorm.Run(bluestorm.CloseOnSignal(os.Interrupt), servers)
}
