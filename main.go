package main

import (
	"log"
	"net"
)

func main() {
	s := newServer()

	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatalf("Unable to strat server: %s", err.Error())
	}

	defer listener.Close()

	log.Printf("Started server on port 8080")

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Printf("Unable to accept connection: %s", err.Error())
			continue
		}
		go s.newClient(conn)
	}
}
