package main

import (
	"bufio"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gorilla/websocket"
)

type MockResponseWriter struct{}

func (mrw MockResponseWriter) Header() http.Header {
	header := http.Header{}
	// header.Add("Connection", "upgrade")
	return header
}

func (mrw MockResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (mrw MockResponseWriter) WriteHeader(statusCode int) {}

func (mrw MockResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return nil, nil, nil
}

func TestFindClients(t *testing.T) {
	con := websocket.Conn{}
	cli := client{"TestName", &con}
	clients = []client{cli}

	foundClient, _ := findClient(&con)
	assert.NotNil(t, foundClient, "Client couldn't be find")
}

func TestFindClientsIfNoClient(t *testing.T) {
	con := websocket.Conn{}
	cli := client{"TestName", &con}
	clients = []client{cli}

	_, err := findClient(nil)
	assert.Error(t, err, "Client shoudn't be find and error should be thrown")
}

// func TestClientIsAddedAfterInitiation(t *testing.T) {
// 	header := http.Header{}
// 	header.Add("Connection", "upgrade")
// 	header.Add("Upgrade", "websocket")
// 	header.Add("Sec-Websocket-Version", "13")
// 	header.Add("Sec-WebSocket-Key", "test")
// 	_, clients := initConnection(MockResponseWriter{}, &http.Request{Header: header, Method: "GET"})
// 	if len(clients) == 0 {
// 		t.Errorf("Client wasnt't added to the slice, clients slice: %v.", clients)
// 	}
// }

// func TestClientIsNotAddedAfterInitiation(t *testing.T) {

// }
