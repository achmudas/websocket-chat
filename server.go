package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/marstr/randname"
)

type client struct {
	id         string
	connection *websocket.Conn
}

var clients = []client{}

type conUpgrader interface {
	upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error)
}

type connectionUpgrader struct{}

func (conUp connectionUpgrader) upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{}
	con, err := upgrader.Upgrade(w, r, nil)
	return con, err
}

func initConnection(w http.ResponseWriter, r *http.Request, cu conUpgrader) (*websocket.Conn, []client) {
	con, err := cu.upgrade(w, r)
	if err != nil {
		log.Fatal("Error when innitiating connection: ", err)
	}

	clientName := randname.Generate()
	fmt.Printf("Connected %s\n", clientName)
	clients = append(clients, client{clientName, con})
	return con, clients
}

func receive(w http.ResponseWriter, r *http.Request) {
	conUpgrader := connectionUpgrader{}
	con, clients := initConnection(w, r, conUpgrader)
	defer con.Close()

	for {
		client, err := findClient(con)
		if err != nil {
			log.Printf("Error when connecting client: %v", err)
		}

		mt, msg, err := con.ReadMessage()
		if err != nil {
			if c, k := err.(*websocket.CloseError); k {
				if c.Code == 1000 {
					fmt.Printf("%s disconnected\n", client.id)
					break
				}
			} else {
				log.Printf("Error when reading message: %v", err)
			}
		}

		fmt.Printf("Received from %s: %s\n", client.id, msg)

		// #FIXME use go routines?
		for _, cli := range clients {
			msgToSent := []byte{}
			if cli.connection == con {
				msgToSent = append([]byte(client.id+": "), msg...)
			} else {
				msgToSent = append([]byte("\n"+client.id+": "), msg...)
			}
			err = cli.connection.WriteMessage(mt, msgToSent)
			if err != nil {
				log.Printf("Error when sending message: %v", err)
			}
		}
	}
}

func findClient(con *websocket.Conn) (*client, error) {
	for _, cli := range clients {
		if cli.connection == con {
			return &cli, nil
		}
	}
	return nil, errors.New("Client couldn't be found")
}

func main() {
	address := "localhost:8080"
	http.HandleFunc("/", receive)
	log.Printf("Server listening on: %v", address)
	log.Fatal(http.ListenAndServe(address, nil)) //#FIXME use flag to let user define address
}
