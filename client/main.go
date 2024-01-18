package main

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

const address = "localhost:8080"

func main() {
	var (
		wsUrl        = url.URL{Scheme: "ws", Host: address, Path: "/"}
		conn, _, err = websocket.DefaultDialer.Dial(wsUrl.String(), nil)
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	log.Println("Connected to", address)

}
