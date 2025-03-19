package server

import (
	"fmt"
	"log"
	"net"
	"os"
)

type Redis struct {
	host     string
	port     int
	listener net.Listener
}

func NewRedis(host string, port int) *Redis {
	return &Redis{
		host,
		port,
		nil,
	}
}

func (r *Redis) Run() error {
	address := fmt.Sprintf("%s:%d", r.host, r.port)

	var err error
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	r.listener = listener

	connection, err := r.listener.Accept()
	if err != nil {
		log.Printf("Error accepting connection: %v", err)
		os.Exit(1)
	}

	for {
		buffer := make([]byte, 1024)
		_, err := connection.Read(buffer)
		if err != nil {
			log.Printf("Error reading from connection: %v", err)
			continue
		}

		_, err = connection.Write([]byte("+PONG\r\n"))
		if err != nil {
			log.Printf("Error writing to connection: %v", err)
			continue
		}
	}
}
