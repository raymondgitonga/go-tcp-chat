package main

import (
	"log"
	"net"
)

type server struct {
	rooms    map[string]*room
	commands chan command
}

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}
func (s *server) newClient(conn net.Conn) {
	log.Printf("new client has connected %s ", conn.RemoteAddr().String())

	c := &client{
		conn:     conn,
		nick:     "anon",
		commands: s.commands,
	}

	c.readInput()
}
