package server

import (
	"fmt"
	"log"
	"net"

	"github.com/codecrafters-io/redis-starter-go/internal/command"
	"github.com/codecrafters-io/redis-starter-go/pkg/protocol"
)

// Redis represents a Redis server instance.
// It handles TCP connections and implements basic Redis functionality.
type Redis struct {
	host     string
	port     int
	listener net.Listener
}

// NewRedis creates a new Redis server instance.
//
// It initializes a Redis server with the specified host and port.
//
// Parameters:
//   - host: The host address to listen on (e.g., "0.0.0.0")
//   - port: The port number to listen on (e.g., 6379)
//
// Returns a new Redis server instance.
func NewRedis(host string, port int) *Redis {
	return &Redis{
		host,
		port,
		nil,
	}
}

// Redis represents a Redis server instance.
// It handles TCP connections and implements basic Redis functionality.
//
// The server listens on the specified host and port, accepting client connections
// and processing Redis protocol commands.
//
// Example usage:
//
//	redis := server.NewRedis("0.0.0.0", 6379)
//	redis.Run()
func (r *Redis) Run() error {
	address := fmt.Sprintf("%s:%d", r.host, r.port)

	var err error
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Printf("Error listening on %s: %v", address, err)
		return err
	}
	r.listener = listener

	for {
		connection, err := r.listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go r.handleConnection(connection)
	}
}

func (r *Redis) handleConnection(connection net.Conn) {
	defer connection.Close()

	for {
		buffer := make([]byte, 1024)
		_, err := connection.Read(buffer)

		if err != nil {
			log.Printf("Error reading from connection: %v", err)
			return
		}

		query := protocol.NewQuery()
		err = query.Deserialize(buffer)
		if err != nil {
			log.Printf("Error deserializing query: %v", err)
			return
		}

		result := command.Processor(query)

		_, err = connection.Write(result.Serialize())
		if err != nil {
			log.Printf("Error writing to connection: %v", err)
			return
		}
	}
}
