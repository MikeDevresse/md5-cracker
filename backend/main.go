package main

import (
	"fmt"
	"github.com/MikeDevresse/md5-cracker/pkg/websocket"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
)

// main init the websocket and starts the webserver
func main() {
	port := 80
	initWebsocket()
	log.Println("main.go", "Running webserver on port", port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		log.Fatal("An error occurred while trying to start the webserver: \"", err, "\"")
	}
}

// initWebsocket initialize the websocket route and the server
func initWebsocket() {
	server := websocket.NewServer(redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	}))
	go server.Start()

	clientCount := 0
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Upgrade(w, r)
		if err != nil {
			fmt.Fprintf(w, "%+V\n", err)
			return
		}

		// Read the first message in order to identify the client
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Fprintf(w, "%+V\n", err)
			return
		} else if string(p) == "slave" {
			slave := websocket.Client{
				Identifier: clientCount,
				Conn:       conn,
				IsSlave:    true,
				Server:     server,
			}
			server.AddSlave(&slave)
			go slave.Read()
		} else {
			client := websocket.Client{
				Identifier: clientCount,
				Conn:       conn,
				IsSlave:    false,
				Server:     server,
			}
			server.AddClient(&client)
			server.PrintConfiguration(&client)
			go client.Read()
		}
		clientCount++
	})
}
