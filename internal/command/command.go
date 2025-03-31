package command

import (
	"strconv"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/store"
	"github.com/codecrafters-io/redis-starter-go/pkg/protocol"
)

const (
	PING = "PING"
	ECHO = "ECHO"
	OK   = "OK"
	ERR  = "ERR"
	SET  = "SET"
	GET  = "GET"
	DEL  = "DEL"
)

var ALLOWED_COMANDS = map[string]struct{}{
	"PING": {},
	"ECHO": {},
	"SET":  {},
	"GET":  {},
	"DEL":  {},
}

// Processor processes Redis commands and returns appropriate responses.
//
// This function takes a Query object containing a Redis command and its arguments,
// validates the command, and executes the appropriate handler to generate a response.
// It supports basic Redis commands like PING, ECHO, SET, GET, and DEL.
//
// Parameters:
//   - query: A pointer to a protocol.Query object containing the command and arguments
//
// Returns a protocol.Result object containing the formatted response.
//
// Example:
//
//	query := &protocol.Query{
//	    Command: "PING",
//	    Args:    []string{},
//	}
//	result := Processor(query) // Returns a SimpleString "PONG" response
func Processor(query *protocol.Query) *protocol.Result {
	if _, ok := ALLOWED_COMANDS[query.Command]; !ok {
		return protocol.NewResult(protocol.ErrorType, "ERR unknown command")
	}

	switch query.Command {
	case PING:
		return protocol.NewResult(protocol.SimpleStringType, "PONG")
	case ECHO:
		if len(query.Args) == 0 {
			return protocol.NewResult(protocol.ErrorType, "ERR wrong number of arguments for 'ECHO' command")
		}
		return protocol.NewResult(protocol.BulkStringType, query.Args[0])
	case SET:
		return runSetCommand(query)
	case GET:
		return runGetCommand(query)
	default:
		return protocol.NewResult(protocol.ErrorType, "ERR unknown command")
	}
}

func runSetCommand(query *protocol.Query) *protocol.Result {
	if len(query.Args) < 2 {
		return protocol.NewResult(protocol.ErrorType, "ERR wrong number of arguments for 'SET' command")
	}

	key := query.Args[0]
	value := query.Args[1]

	var expiryTime time.Duration = 0

	for i := 2; i < len(query.Args); i++ {
		option := query.Args[i]

		if option == "EX" || option == "PX" {
			if i+1 >= len(query.Args) {
				return protocol.NewResult(protocol.ErrorType, "ERR syntax error")
			}

			expStr := query.Args[i+1]
			expVal, err := strconv.Atoi(expStr)
			if err != nil || expVal <= 0 {
				return protocol.NewResult(protocol.ErrorType, "ERR value is not an integer or out of range")
			}

			if option == "EX" {
				expiryTime = time.Duration(expVal) * time.Second
			} else {
				expiryTime = time.Duration(expVal) * time.Millisecond
			}

			i++
		}
	}

	store.Set(key, value, expiryTime)
	return protocol.NewResult(protocol.SimpleStringType, OK)
}

func runGetCommand(query *protocol.Query) *protocol.Result {
	if len(query.Args) != 1 {
		return protocol.NewResult(protocol.ErrorType, "ERR wrong number of arguments for 'GET' command")
	}
	key := query.Args[0]
	value, exists := store.Get(key)
	if !exists {
		return protocol.NewResult(protocol.BulkStringType, "")
	}
	return protocol.NewResult(protocol.BulkStringType, value)
}
