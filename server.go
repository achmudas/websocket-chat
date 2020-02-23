package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/marstr/randname"
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

	clientName := randname.Generate()

	fmt.Printf("Connected %s\n", clientName)

	clients = append(clients, Client{clientName, con})

	for {

		var client Client
		for _, cli := range clients {
			if cli.connection == con {
				client = cli
			}
		}

		mt, msg, err := con.ReadMessage()
		if err != nil {
			if c, k := err.(*websocket.CloseError); k {
				if c.Code == 1000 {
					fmt.Printf("%s disconnected\n", client.id)
					break
				}
			} else {
				log.Printf("Error when reading message: ", err)
			}
		}

		fmt.Printf("MT: %d\n", mt)

		// sendingClient := getSendingClient(con)

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
				log.Printf("Error when sending message: ", err)
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
