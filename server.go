package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
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

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			s.nick(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.listRooms(cmd.client, cmd.args)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client, cmd.args)
		}
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

func (s *server) join(c *client, args []string) {
	c.nick = args[1]
	c.msg(fmt.Sprintf("Your name is now %s", c.nick))

}

func (s *server) nick(c *client, args []string) {
	// Get room name
	roomName := args[1]

	// Check if room exists
	r, ok := s.rooms[roomName]

	//If not add to rooms property
	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}

		s.rooms[roomName] = r
	}

	//add current client to list of members
	r.members[c.conn.RemoteAddr()] = c

	// Remove user from other room there in
	s.quitCurrentRoom(c)

	//assign new room to client
	c.room = r
	r.broadcast(c, fmt.Sprintf("%s, just joined the room", c.nick))
	c.msg(fmt.Sprintf("Welcome to room %s", r.name))
}

func (s *server) listRooms(c *client, args []string) {
	var rooms []string

	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	c.msg(fmt.Sprintf("available rooms %s", strings.Join(rooms, " ,")))
}

func (s *server) msg(c *client, args []string) {
	if c.room != nil {
		c.err(errors.New("you must join room first"))
		return
	}
	c.room.broadcast(c, c.nick+": "+strings.Join(args[1:len(args)], " "))
}

func (s *server) quit(c *client, args []string) {
	log.Printf("client has disconnected: %s" + c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c)

	c.conn.Close()
}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		delete(c.room.members, c.conn.RemoteAddr())
		c.room.broadcast(c, fmt.Sprintf("%s, has left the room", c.nick))
	}
}
