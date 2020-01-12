package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"

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

	reader := bufio.NewReader(os.Stdin)

	go func() {
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Printf("Error when reading message: ", err)
			}
			fmt.Printf("%s> ", msg)
		}
	}()

	fmt.Printf("> ")

	for {
		bytes, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error when reading std input: ", err)
		}
		c.WriteMessage(websocket.TextMessage, bytes)
	}
}
