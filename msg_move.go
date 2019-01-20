package main

type MsgMove struct {
	ID  string `json:"id"`
	X   int    `json:"x"`
	Y   int    `json:"y"`
	Dir int    `json:"dir"`
}

func (msg *MsgMove) Run(receiver Receiver, player *Player) {
	player.nextDir = msg.Dir
}
