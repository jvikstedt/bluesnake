package main

import "log"

type MsgPing struct {
}

func (m *MsgPing) Run(receiver Receiver, msg *Msg) {
	log.Printf("Run msg ping %v %v\n", receiver, msg)

	// lobby, ok := receiver.(*Lobby)
	// if !ok {
	// 	log.Printf("Could not convert to lobby...%v\n", receiver)
	// 	return
	// }

	if err := msg.Player.WriteMsg(&MsgPing{}); err != nil {
		log.Println(err)
	}
}
