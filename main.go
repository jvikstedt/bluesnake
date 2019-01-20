package main

import (
	"log"
	"os"

	"github.com/jvikstedt/bluestorm"
	"github.com/jvikstedt/bluestorm/hub"
	"github.com/jvikstedt/bluestorm/network"
	"github.com/jvikstedt/bluestorm/network/json"
)

const LobbyID = "lobby"

func proxy(agent *network.Agent, i interface{}) {
	user, err := hub.DefaultManager().GetUser(hub.UserID(agent.ID()))
	if err != nil {
		log.Println(err)
		return
	}

	player, ok := user.(*Player)
	if !ok {
		log.Println("Could not convert to Player...")
		return
	}

	room := player.GetRoom()

	receiver, ok := room.(Receiver)
	if !ok {
		log.Printf("Could not convert to cmd receiver... %v\n", room)
		return
	}

	command, ok := i.(Command)
	if !ok {
		log.Printf("Could not convert to command...%v\n", i)
		return
	}

	receiver.AddMsg(&Msg{
		Player:  player,
		Command: command,
	})
}

func onConnect(agent *network.Agent) {
	player := &Player{
		BaseUser: hub.NewBaseUser(hub.UserID(agent.ID())),
		agent:    agent,
	}

	if err := hub.DefaultManager().AddUser(player, LobbyID); err != nil {
		agent.Conn().Close()
		return
	}
}

func onDisconnect(agent *network.Agent) {
	if err := hub.DefaultManager().RemoveUser(hub.UserID(agent.ID())); err != nil {
		log.Println(err)
		return
	}
}

func main() {
	processor := json.NewProcessor()
	processor.Register(&MsgPing{}, proxy)

	lobby := NewLobby(LobbyID)
	hub.DefaultManager().AddRoom(lobby)

	servers := bluestorm.Servers{
		&network.WSServer{
			Addr:         ":8081",
			Processor:    processor,
			OnConnect:    onConnect,
			OnDisconnect: onDisconnect,
		},
	}

	go lobby.Run()

	bluestorm.Run(bluestorm.CloseOnSignal(os.Interrupt), servers)
}
