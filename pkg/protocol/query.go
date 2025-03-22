package protocol

import "fmt"

const (
	SimpleString = '+'
	Error        = '-'
	Integer      = ':'
	BulkString   = '$'
	Array        = '*'
)

// Query represents a Redis protocol query with a command and its arguments.
//
// The Query struct is used to parse and represent Redis protocol commands
// received from clients. It contains the main command (like "GET", "SET", "PING")
// and any arguments that follow the command.
//
// Example:
//
//	query := &Query{
//	    Command: "SET",
//	    Args:    []string{"key", "value"},
//	}
type Query struct {
	Command string
	Args    []string
}

// NewQuery creates a new Query object.
//
// This function initializes a new Query instance with empty command and arguments.
//
// Returns a pointer to the newly created Query object.
func NewQuery() *Query {
	return &Query{}
}

// Deserialize parses Redis protocol data into a Query object.
//
// It handles different Redis protocol types (SimpleString, Error, Integer, BulkString, Array)
// and extracts the command and arguments from the raw byte data.
//
// Parameters:
//   - data: Raw byte data in Redis protocol format
//
// Returns an error if the data is malformed or cannot be parsed.
func (q *Query) Deserialize(data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("empty data received")
	}

	switch data[0] {
	case SimpleString:
		q.Command = string(data[1 : len(data)-2])
		q.Args = []string{}
		return nil

	case Error:
		return fmt.Errorf("received error: %s", string(data[1:len(data)-2]))

	case Integer:
		q.Command = string(data[1 : len(data)-2])
		q.Args = []string{}
		return nil

	case BulkString:
		firstCRLF := -1
		for i := 1; i < len(data)-1; i++ {
			if data[i] == '\r' && data[i+1] == '\n' {
				firstCRLF = i
				break
			}
		}

		if firstCRLF == -1 {
			return fmt.Errorf("malformed bulk string: missing length delimiter")
		}

		length := 0
		for i := 1; i < firstCRLF; i++ {
			if data[i] < '0' || data[i] > '9' {
				return fmt.Errorf("malformed bulk string: invalid length")
			}
			length = length*10 + int(data[i]-'0')
		}

		if length == -1 {
			q.Command = ""
			q.Args = []string{}
			return nil
		}

		start := firstCRLF + 2
		if start+length+2 > len(data) {
			return fmt.Errorf("malformed bulk string: insufficient data")
		}

		q.Command = string(data[start : start+length])
		q.Args = []string{}
		return nil

	case Array:
		firstCRLF := -1
		for i := 1; i < len(data)-1; i++ {
			if data[i] == '\r' && data[i+1] == '\n' {
				firstCRLF = i
				break
			}
		}

		if firstCRLF == -1 {
			return fmt.Errorf("malformed array: missing count delimiter")
		}

		count := 0
		for i := 1; i < firstCRLF; i++ {
			if data[i] < '0' || data[i] > '9' {
				return fmt.Errorf("malformed array: invalid count")
			}
			count = count*10 + int(data[i]-'0')
		}

		if count == -1 {
			q.Command = ""
			q.Args = []string{}
			return nil
		}

		if count == 0 {
			q.Command = ""
			q.Args = []string{}
			return nil
		}

		elements := make([]string, 0, count)
		pos := firstCRLF + 2

		for i := 0; i < count; i++ {
			if pos >= len(data) {
				return fmt.Errorf("malformed array: insufficient data")
			}

			if data[pos] != '$' {
				return fmt.Errorf("malformed array: expected bulk string")
			}

			nextCRLF := -1
			for j := pos + 1; j < len(data)-1; j++ {
				if data[j] == '\r' && data[j+1] == '\n' {
					nextCRLF = j
					break
				}
			}

			if nextCRLF == -1 {
				return fmt.Errorf("malformed array: missing length delimiter")
			}

			length := 0
			for j := pos + 1; j < nextCRLF; j++ {
				if data[j] < '0' || data[j] > '9' {
					return fmt.Errorf("malformed array: invalid length")
				}
				length = length*10 + int(data[j]-'0')
			}

			start := nextCRLF + 2
			if start+length+2 > len(data) {
				return fmt.Errorf("malformed array: insufficient data")
			}

			elements = append(elements, string(data[start:start+length]))
			pos = start + length + 2
		}

		if len(elements) > 0 {
			q.Command = elements[0]
			q.Args = elements[1:]
		} else {
			q.Command = ""
			q.Args = []string{}
		}
		return nil

	default:
		return fmt.Errorf("unknown RESP data type: %c", data[0])
	}
}
