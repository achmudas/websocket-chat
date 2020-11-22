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

type connectionUpgraderMock struct{}

func (conUp connectionUpgraderMock) upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	return upgradeMock()
}

var upgradeMock func() (*websocket.Conn, error)

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

func TestClientIsAddedAfterInitiation(t *testing.T) {
	upMock := connectionUpgraderMock{}
	upgradeMock = func() (*websocket.Conn, error) {
		return nil, nil
	}
	_, clients := initConnection(MockResponseWriter{}, &http.Request{Header: nil, Method: "GET"}, upMock)
	assert.NotEmpty(t, clients, "Client wasnt't added to the slice")
}

// func TestClientIsNotAddedAfterInitiation(t *testing.T) {

// }
