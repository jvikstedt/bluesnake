package main

import (
	"log"
	"os"

	"github.com/jvikstedt/bluestorm"
	"github.com/jvikstedt/bluestorm/hub"
	"github.com/jvikstedt/bluestorm/network"
	"github.com/jvikstedt/bluestorm/network/json"
)

type Player struct {
	*hub.BaseUser

	agent *network.Agent
}

func (p *Player) WriteMsg(i interface{}) error {
	return p.agent.WriteMsg(i)
}

type Msg struct {
	Message string `json:"name"`
}

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

	lobby, ok := room.(*Lobby)
	if !ok {
		log.Println("Could not convert to lobby...")
		return
	}

	lobby.NewMsg(player, i)
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
	processor.Register(&Msg{}, proxy)

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
