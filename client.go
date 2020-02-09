package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/achmudas/websocket-chat/commands"
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
			// fmt.Println("MT %d", msgType.Code)
			fmt.Printf("%s> ", msg)
		}
	}()

	fmt.Printf("> ")

	for {
		peakByte, err := reader.Peek(1)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error when peeking byte std input: ", err)
		}

		if peakByte[0] == byte(47) {
			commandBytes, err := reader.ReadBytes('\n')
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to disconnect from server: ", err)
			}

			command, err := commands.Create(string(commandBytes[1:]))
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to execute command: ", err)
				fmt.Printf("> ")
				// return
			}
			if command != nil {
				command.Execute(c)
			}
			// break #FIXME ?
		}

		bytes, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error when reading std input: ", err)
		}
		c.WriteMessage(websocket.TextMessage, bytes)
	}
}
