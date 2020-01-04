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
	return true
}

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
	for {
		mt, msg, err := con.ReadMessage()
		if err != nil {
			log.Printf("Error when reading message: ", err)
		}

		fmt.Printf("Received: %s", msg)

		err = con.WriteMessage(mt, msg)
		if err != nil {
			log.Printf("Error when sending echo message: ", err)
		}
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
