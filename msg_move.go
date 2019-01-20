package main

import "log"

type MsgMove struct {
	ID  string `json:"id"`
	X   int    `json:"x"`
	Y   int    `json:"y"`
	Dir int    `json:"dir"`
}

func (msg *MsgMove) Run(receiver Receiver, player *Player) {
	log.Printf("move %d\n", msg.Dir)
	player.nextDir = msg.Dir
}
