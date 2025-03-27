package commands

import (
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"time"
)

type MockConn struct {
	*websocket.Conn
}

func (m *MockConn) WriteControl(messageType int, data []byte, deadline time.Time) error {
	return nil
}

func TestCreateDisconnectCommand(t *testing.T) {
	cmd, err := Create("quit")
	assert.NoError(t, err, "Creating 'quit' command should not produce an error")
	assert.IsType(t, Disconnect{}, cmd, "Command should be of type Disconnect")
}

func TestCreateUnknownCommand(t *testing.T) {
	cmd, err := Create("unknown")
	assert.Error(t, err, "Creating an unknown command should produce an error")
	assert.Nil(t, cmd, "Command should be nil for unknown command")
}

func TestDisconnectCommandExecution(t *testing.T) {
	conn := &MockConn{}
	cmd := Disconnect{}
	quit, err := cmd.Execute(conn)
	assert.NoError(t, err, "Executing Disconnect command should not produce an error")
	assert.True(t, quit, "Executing Disconnect command should return true for quit")
}
