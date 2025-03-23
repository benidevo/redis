package server

import (
	"net"
	"strconv"
	"testing"
	"time"
)

func TestRedisServer(t *testing.T) {

	t.Run("basic connection and PING", func(t *testing.T) {
		_, serverAddress := setupTestRedisServer(t)

		conn := connectToRedis(t, serverAddress)
		defer conn.Close()

		pingCmd := []byte("*1\r\n$4\r\nPING\r\n")
		_, err := conn.Write(pingCmd)
		if err != nil {
			t.Fatalf("Failed to write to connection: %v", err)
		}

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			t.Fatalf("Failed to read from connection: %v", err)
		}

		response := buffer[:n]
		expectedResponse := []byte("+PONG\r\n")

		if string(response) != string(expectedResponse) {
			t.Fatalf("Unexpected response: got %q, want %q", response, expectedResponse)
		}
	})
	t.Run("ECHO command", func(t *testing.T) {
		_, serverAddress := setupTestRedisServer(t)

		conn := connectToRedis(t, serverAddress)
		defer conn.Close()

		echoCmd := []byte("*2\r\n$4\r\nECHO\r\n$11\r\nHello Redis\r\n")
		_, err := conn.Write(echoCmd)
		if err != nil {
			t.Fatalf("Failed to write to connection: %v", err)
		}

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			t.Fatalf("Failed to read from connection: %v", err)
		}

		response := buffer[:n]
		expectedResponse := []byte("$11\r\nHello Redis\r\n")

		if string(response) != string(expectedResponse) {
			t.Fatalf("Unexpected response: got %q, want %q", response, expectedResponse)
		}
	})
	t.Run("Invalid command", func(t *testing.T) {
		_, serverAddress := setupTestRedisServer(t)

		conn := connectToRedis(t, serverAddress)
		defer conn.Close()

		invalidCmd := []byte("*1\r\n$7\r\nINVALID\r\n")
		_, err := conn.Write(invalidCmd)
		if err != nil {
			t.Fatalf("Failed to write to connection: %v", err)
		}

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			t.Fatalf("Failed to read from connection: %v", err)
		}

		response := buffer[:n]
		expectedResponse := []byte("-ERR unknown command\r\n")

		if string(response) != string(expectedResponse) {
			t.Fatalf("Unexpected response: got %q, want %q", response, expectedResponse)
		}
	})
}

func setupTestRedisServer(t *testing.T) (*Redis, string) {
	t.Helper()
	redis := NewRedis("localhost", 0)
	go func() {
		err := redis.Run()
		if err != nil {
			t.Errorf("Failed to start server: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)
	addr := redis.listener.Addr().(*net.TCPAddr)
	serverAddress := "localhost:" + strconv.Itoa(addr.Port)
	return redis, serverAddress
}

func connectToRedis(t *testing.T, serverAddress string) net.Conn {
	t.Helper()

	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	return conn
}

type MockConnection struct {
	ReadData  []byte
	WriteData []byte
	Closed    bool
}

func (m *MockConnection) Read(b []byte) (n int, err error) {
	copy(b, m.ReadData)
	return len(m.ReadData), nil
}

func (m *MockConnection) Write(b []byte) (n int, err error) {
	m.WriteData = make([]byte, len(b))
	copy(m.WriteData, b)
	return len(b), nil
}

func (m *MockConnection) Close() error {
	m.Closed = true
	return nil
}

func (m *MockConnection) LocalAddr() net.Addr {
	return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0}
}

func (m *MockConnection) RemoteAddr() net.Addr {
	return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0}
}

func (m *MockConnection) SetDeadline(t time.Time) error {
	return nil
}

func (m *MockConnection) SetReadDeadline(t time.Time) error {
	return nil
}

func (m *MockConnection) SetWriteDeadline(t time.Time) error {
	return nil
}

func TestHandleConnection(t *testing.T) {
	redis := NewRedis("localhost", 6379)

	mockConn := &MockConnection{
		ReadData: []byte("*1\r\n$4\r\nPING\r\n"),
	}

	go redis.handleConnection(mockConn)

	time.Sleep(100 * time.Millisecond)

	expectedResponse := []byte("+PONG\r\n")
	if string(mockConn.WriteData) != string(expectedResponse) {
		t.Fatalf("Unexpected response: got %q, want %q", mockConn.WriteData, expectedResponse)
	}
}
