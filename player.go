package main

import (
	"github.com/jvikstedt/bluestorm/hub"
	"github.com/jvikstedt/bluestorm/network"
)

type Player struct {
	*hub.BaseUser

	agent   *network.Agent
	x       int
	y       int
	dir     int
	nextDir int
}

func (p *Player) WriteMsg(i interface{}) error {
	return p.agent.WriteMsg(i)
}
