package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func handleError(w http.ResponseWriter, r *http.Request, status int, reason error) {
	// #FIXME finish this one
}

func checkOrigin(r *http.Request) bool {
	//#FIXME ignored for now
	return false
}

func initConnection() *websocket.Upgrader {
	upgrader := websocket.Upgrader{
		10,
		1024,
		1024,
		nil,
		nil,
		handleError,
		checkOrigin,
		false,
	}
	return &upgrader
}

func receive(w http.ResponseWriter, r *http.Request) {
	up := initConnection()
	con, err := up.Upgrade(w, r, nil) //#FIXME Upgrade is deprecated
	if err != nil {
		log.Printf("Error when innitiating connection: %v", err)
	}
	// defer con.Close()
	for {
		_, msg, err := con.ReadMessage()
		if err != nil {
			log.Printf("Error when reading message: %v", err)
		}
		fmt.Printf("Received: %s", msg)
	}
}

func main() {
	address := "localhost:8080"
	http.HandleFunc("/", receive)
	log.Printf("Server listening on: %v", address)
	log.Fatal(http.ListenAndServe(address, nil)) //#FIXME use flag to let user define address
	// reader := bufio.NewReader(os.Stdin)
	// for {
	// 	fmt.Print("Message: ")
	// 	input, err := reader.ReadString('\n')
	// 	if err != nil {
	// 		log.Panicf("Error when reading input: %v", err)
	// 	}
	// 	fmt.Printf("User ententered: %v", input)
	// }
}
