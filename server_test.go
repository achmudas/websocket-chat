package main

import (
	"testing"

	"github.com/gorilla/websocket"
)

func TestFindClients(t *testing.T) {
	con := websocket.Conn{}
	cli := client{"TestName", &con}
	clients = []client{cli}

	foundClient, _ := findClient(&con)
	if foundClient == nil {
		t.Errorf("Client couldn't be find, got: %v, want: %v.", foundClient, cli)
	}

}
