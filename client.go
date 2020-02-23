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
	c := connect()
	defer c.Close()

	reader := bufio.NewReader(os.Stdin)

	go waitAndReadMessages(c)

	fmt.Printf("> ")

	for {
		peakByte, err := reader.Peek(1)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error when peeking byte std input: ", err)
		}

		// #FIXME move to separate function?
		if peakByte[0] == byte(47) {
			commandBytes, err := reader.ReadBytes('\n')
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to disconnect from server: ", err)
			}

			command, err := commands.Create(string(commandBytes[1:]))
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to find command: ", err)
			}
			if command != nil {
				quit, err := command.Execute(c)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Failed to execute command: ", err)
				}

				if quit {
					break
				}

			}
			fmt.Printf("> ")
		}

		bytes, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error when reading std input: ", err)
		}
		c.WriteMessage(websocket.TextMessage, bytes)
	}
}

func connect() (conn *websocket.Conn) {
	u := url.URL{Scheme: "ws", Host: "localhost:8080"}

	dialer := websocket.DefaultDialer
	c, _, err := dialer.Dial(u.String(), nil)

	if err != nil {
		log.Fatal("Failed to connect ", err)
	}
	return c
}

func waitAndReadMessages(c *websocket.Conn) {
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Printf("Error when reading message: ", err)
		}
		fmt.Printf("%s> ", msg)
	}
}
