package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/achmudas/websocket-chat/client/commands"
	"github.com/gorilla/websocket"
)

func main() {
	c := connect()

	reader := bufio.NewReader(os.Stdin)
	go waitAndReadMessages(c)
	fmt.Printf("> ")

	for {
		peekByte, err := reader.Peek(1)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error when peeking byte std input: ", err)
		}

		if quit := executeFunctionalCommand(peekByte, reader, c); quit {
			c.Close()
			break
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
			log.Printf("Error when reading message: %v ", err)
		}
		fmt.Printf("%s> ", msg)
	}
}

func executeFunctionalCommand(peekByte []byte, reader *bufio.Reader, c *websocket.Conn) bool {
	fmt.Println(peekByte[0])
	if peekByte[0] == byte(47) {
		commandBytes, err := reader.ReadBytes('\n')
		fmt.Println(commandBytes)
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
			return quit
		}
		fmt.Printf("> ")
	}
	return false
}
