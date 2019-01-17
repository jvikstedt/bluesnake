package main

import (
	"log"
	"os"

	"github.com/jvikstedt/bluestorm"
	"github.com/jvikstedt/bluestorm/network"
	"github.com/jvikstedt/bluestorm/network/json"
)

func onConnect(a *network.Agent) {
	log.Printf("Agent connectedted %s\n", a.ID)
}

func onDisconnect(a *network.Agent) {
	log.Printf("Agent disconnected %s\n", a.ID)
}

func main() {
	processor := json.NewProcessor()
	// processor.Register(&Greet{}, test)

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
