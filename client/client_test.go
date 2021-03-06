package main

import (
	"bufio"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestCommandIsFound(t *testing.T) {
	// reader := bufio.Reader{}
	// executeFunctionalCommand([]byte{byte(47)}, &reader, &websocket.Conn{})
}

func TestIsNotFunctionalCommandAndNotDisconnect(t *testing.T) {
	reader := bufio.Reader{}

	quit := executeFunctionalCommand([]byte{byte(48)}, &reader, &websocket.Conn{})
	assert.False(t, quit, "It's not a functional command and client shouldn't disconnect")
}

func TestIsDisconnectCommand(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader("/quit\n"))
	quit := executeFunctionalCommand([]byte{byte(47)}, reader, &websocket.Conn{})
	assert.True(t, quit, "It's a functional 'quit' command and client should disconnect")
}

func TestCommandIsNotFound(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader("/amazing\n"))
	quit := executeFunctionalCommand([]byte{byte(47)}, reader, &websocket.Conn{})
	assert.False(t, quit, "It's a functional non existing command")
}
