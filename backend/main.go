package main

import (
	"fmt"
	"github.com/MikeDevresse/md5-cracker/pkg/websocket"
	"log"
	"net/http"
)

func main() {
	setupRoutes()
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("An error occurred while trying to start the webserver: \"", err, "\"")
	}

}

func setupRoutes() {
	server := websocket.NewServer()
	go server.Start()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "MD5 Cracker")
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Upgrade(w, r)
		if err != nil {
			fmt.Fprintf(w, "%+V\n", err)
		}

		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error while reading welcome message: ", p)
			return
		}

		if string(p) == "slave" {
			client := &websocket.Client{
				Conn:    conn,
				Pool:    server.SlavePool,
				Server:  server,
				IsSlave: true,
			}

			server.SlavePool.Register <- client
			client.Read()
		} else {
			client := &websocket.Client{
				Conn:    conn,
				Pool:    server.ClientPool,
				Server:  server,
				IsSlave: false,
			}

			server.ClientPool.Register <- client
			client.Read()
		}
	})
}
