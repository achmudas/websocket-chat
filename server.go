package main

import (
	"fmt"
	"log"
	"net/http"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/gorilla/websocket"
)

type Client struct {
	id         string
	connection *websocket.Conn
}

var clients = []Client{}

func initConnection() *websocket.Upgrader {
	upgrader := websocket.Upgrader{}
	return &upgrader
}

func receive(w http.ResponseWriter, r *http.Request) {
	up := initConnection()
	con, err := up.Upgrade(w, r, nil) //#FIXME Upgrade is deprecated
	if err != nil {
		log.Fatal("Error when innitiating connection: ", err)
	}
	defer con.Close()

	clientName := petname.Name()

	fmt.Printf("Connected %s\n", clientName)

	clients = append(clients, Client{clientName, con})

	for {
		mt, msg, err := con.ReadMessage()
		if err != nil {
			log.Printf("Error when reading message: ", err)
		}

		// sendingClient := getSendingClient(con)

		var client Client
		for _, cli := range clients {
			if cli.connection == con {
				client = cli
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
				log.Printf("Error when sending echo message: ", err)
			}
		}
	}
}

// func getSendingClient(con *websocket.Conn) Client {
// 	var client Client
// 	for _, cli := range clients {
// 		if cli.connection == con {
// 			client = cli
// 		}
// 	}
// 	return client
// }

func main() {
	address := "localhost:8080"
	http.HandleFunc("/", receive)
	log.Printf("Server listening on: %v", address)
	log.Fatal(http.ListenAndServe(address, nil)) //#FIXME use flag to let user define address
}
