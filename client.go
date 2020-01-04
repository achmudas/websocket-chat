package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

func main() {

	u := url.URL{Scheme: "ws", Host: "localhost:8080"}

	dialer := websocket.DefaultDialer
	c, _, err := dialer.Dial(u.String(), nil)

	if err != nil {
		log.Fatal("Failed to connect ", err)
	}

	defer c.Close()

	for {

		c.WriteMessage(websocket.TextMessage, []byte("Hello world!"))

		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Printf("Error when reading message: ", err)
		}
		fmt.Printf("Received response from server: %s", msg)

	}

}
