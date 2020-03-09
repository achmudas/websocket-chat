package commands

import (
	"errors"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// Command is used to declare functional webchat commands, e.g. Disconnect
type Command interface {
	Execute(c *websocket.Conn) (quit bool, err error)
}

type Disconnect struct{}

func (q Disconnect) Execute(c *websocket.Conn) (bool, error) {
	err := c.WriteControl(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Disconnecting"),
		time.Now().Add(time.Second*10))
	return true, err
}

// Create is used to create particular command according passed in user input
func Create(command string) (Command, error) {
	// #FIXME add more options
	switch strings.TrimSpace(command) {
	case "quit":
		return Disconnect{}, nil
	default:
		return nil, errors.New("Command is unrecognized: " + command)
	}
}
