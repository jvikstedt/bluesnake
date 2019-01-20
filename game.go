package main

import (
	"log"
	"time"

	"github.com/jvikstedt/bluestorm/hub"
)

const MaxHeight = 100
const MaxWidth = 100

type Game struct {
	*hub.BaseRoom

	msgs chan *Msg
}

func NewGame(id hub.RoomID) *Game {
	return &Game{
		BaseRoom: hub.NewBaseRoom(id),
		msgs:     make(chan *Msg, 10),
	}
}

func (g *Game) Run() {
	for {
		time.Sleep(time.Second * 1)

		users := g.GetUsers()

	loop:
		for {
			select {
			case msg := <-g.msgs:
				msg.cmd.Run(g, msg.player)
			default:
				break loop
			}
		}

		for _, user := range users {
			player, ok := user.(*Player)
			if !ok {
				log.Printf("Could not convert %v to Player\n", user)
			}

			g.updatePlayerDir(player)
			g.updatePlayerPos(player)
		}

		for _, user := range users {
			player, ok := user.(*Player)
			if !ok {
				log.Printf("Could not convert %v to Player\n", user)
			}

			g.Broadcast(&MsgMove{
				ID:  string(player.ID()),
				X:   player.x,
				Y:   player.y,
				Dir: player.dir,
			})
		}
	}
}

func (g *Game) updatePlayerDir(player *Player) {
	switch player.nextDir {
	case 0:
		if player.dir == 2 {
			return
		}
	case 1:
		if player.dir == 3 {
			return
		}
	case 2:
		if player.dir == 0 {
			return
		}
	case 3:
		if player.dir == 1 {
			return
		}
	default:
		return
	}

	player.dir = player.nextDir
}

func (g *Game) updatePlayerPos(player *Player) {
	newX := player.x
	newY := player.y

	switch player.dir {
	case 0:
		newY--
	case 1:
		newX++
	case 2:
		newY++
	case 3:
		newX--
	default:
		break
	}

	if newX > MaxWidth || newX < 0 || newY > MaxHeight || newY < 0 {
		return
	}

	player.x = newX
	player.y = newY
}

func (g *Game) AddMsg(msg *Msg) {
	g.msgs <- msg
}
