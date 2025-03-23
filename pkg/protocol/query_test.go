package protocol

import (
	"testing"
)

func TestQueryDeserialization(t *testing.T) {
	tests := []struct {
		name           string
		input          []byte
		expectedCmd    string
		expectedArgs   []string
		expectedErrMsg string
	}{
		{
			name:         "Simple PING",
			input:        []byte("*1\r\n$4\r\nPING\r\n"),
			expectedCmd:  "PING",
			expectedArgs: []string{},
		},
		{
			name:         "ECHO with argument",
			input:        []byte("*2\r\n$4\r\nECHO\r\n$5\r\nHello\r\n"),
			expectedCmd:  "ECHO",
			expectedArgs: []string{"Hello"},
		},
		{
			name:           "Malformed command",
			input:          []byte("*x\r\n"),
			expectedErrMsg: "malformed array: invalid count",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := NewQuery()
			err := query.Deserialize(tt.input)

			if tt.expectedErrMsg != "" {
				if err == nil || err.Error() != tt.expectedErrMsg {
					t.Fatalf("Expected error %q, got %v", tt.expectedErrMsg, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if query.Command != tt.expectedCmd {
				t.Errorf("Expected command %q, got %q", tt.expectedCmd, query.Command)
			}

			if len(query.Args) != len(tt.expectedArgs) {
				t.Fatalf("Expected %d args, got %d", len(tt.expectedArgs), len(query.Args))
			}

			for i, arg := range tt.expectedArgs {
				if query.Args[i] != arg {
					t.Errorf("Expected arg %d to be %q, got %q", i, arg, query.Args[i])
				}
			}
		})
	}
}
