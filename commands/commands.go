package commands

import (
	"errors"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// Command is used to declare functional webchat commands, e.g. Disconnect
type Command interface {
	Execute(c *websocket.Conn) error
}

type disconnect struct{}

func (q disconnect) Execute(c *websocket.Conn) error {
	err := c.WriteControl(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Disconnecting"),
		time.Now().Add(time.Second*10))
	return err
}

// Create is used to create particular command according passed in user input
func Create(command string) (Command, error) {
	// #FIXME add more options
	switch strings.TrimSpace(command) {
	case "quit":
		return disconnect{}, nil
	default:
		return nil, errors.New("Command is unrecognized: " + command)
	}
}
