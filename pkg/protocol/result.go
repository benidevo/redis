package protocol

import "fmt"

type ResultType byte

const (
	SimpleStringType ResultType = '+'
	ErrorType        ResultType = '-'
	IntegerType      ResultType = ':'
	BulkStringType   ResultType = '$'
	ArrayType        ResultType = '*'
)

// Result represents a Redis protocol response with a type and value.
//
// The Result struct is used to format and serialize responses according to the Redis
// protocol specification. It contains a type (like SimpleString, Error, Integer, etc.)
// and a value that can be of different types depending on the response type.
//
// The Value field can contain:
// - string: For SimpleString and Error types
// - int: For Integer type
// - string: For BulkString type (serialized with length prefix)
// - []Result: For Array type (nested results)
//
// Example:
//
//	result := &Result{
//	    Type:  SimpleStringType,
//	    Value: "OK",
//	}
//	serialized := result.Serialize() // Returns "+OK\r\n"
type Result struct {
	Type  ResultType
	Value any // This could be string, int, []byte, or []Result
}

// NewResult creates a new Result object.
//
// This function initializes a new Result instance with the specified type and value.
//
// Parameters:
//   - typ: The type of the result (SimpleString, Error, Integer, BulkString, Array)
//   - value: The value of the result (string, int, []byte, []Result)
//
// Returns a pointer to the newly created Result object.
func NewResult(typ ResultType, value any) *Result {
	return &Result{Type: typ, Value: value}
}

// Serialize converts the Result object into a byte slice representing the Redis protocol response.
//
// This method handles different result types (SimpleString, Error, BulkString, Array)
// and formats them according to the Redis protocol specification.
//
// Returns a byte slice containing the serialized result.
func (r *Result) Serialize() []byte {
	switch r.Type {
	case SimpleStringType:
		return []byte("+" + r.Value.(string) + "\r\n")
	case ErrorType:
		return []byte("-" + r.Value.(string) + "\r\n")
	case BulkStringType:
		str := r.Value.(string)
		return fmt.Appendf(nil, "$%d\r\n%s\r\n", len(str), str)
	}
	return nil
}
