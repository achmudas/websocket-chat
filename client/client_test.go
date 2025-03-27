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

type MockConn struct {
	websocket.Conn
}

func (m *MockConn) WriteMessage(messageType int, data []byte) error {
	return nil
}

func (m *MockConn) ReadMessage() (messageType int, p []byte, err error) {
	return websocket.TextMessage, []byte("mock message"), nil
}

func (m *MockConn) ReadMessage() (messageType int, p []byte, err error) {
	return websocket.TextMessage, []byte("mock message"), nil
}

func TestWaitAndReadMessages(t *testing.T) {
	conn := &MockConn{}
	waitAndReadMessages(conn)

	// Simulate reading a message
	messageType, message, err := conn.ReadMessage()
	assert.NoError(t, err, "Reading message should not produce an error")
	assert.Equal(t, websocket.TextMessage, messageType, "Message type should be TextMessage")
	assert.Equal(t, "mock message", string(message), "Message content should match")

	// Simulate reading a message
	messageType, message, err := conn.ReadMessage()
	assert.NoError(t, err, "Reading message should not produce an error")
	assert.Equal(t, websocket.TextMessage, messageType, "Message type should be TextMessage")
	assert.Equal(t, "mock message", string(message), "Message content should match")
}
