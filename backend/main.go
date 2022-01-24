package main

import (
	"fmt"
	"github.com/MikeDevresse/md5-cracker/pkg/websocket"
	"log"
	"net/http"
)

// main call setupRoutes and launch the webserver on port 80
func main() {
	setupRoutes()
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("An error occurred while trying to start the webserver: \"", err, "\"")
	}

}

// setupRoutes set up the websocket route by upgrading http connection to ws one
func setupRoutes() {
	server := websocket.NewServer()
	go server.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Upgrade(w, r)
		if err != nil {
			fmt.Fprintf(w, "%+V\n", err)
		}

		// Check the first message in order to redirect the user to the good pool
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error while reading welcome message: ", p)
			return
		}

		if string(p) == "slave" {
			log.Println("new slave")
			client := &websocket.Client{
				Conn:    conn,
				Pool:    server.SlavePool,
				Server:  server,
				IsSlave: true,
			}

			// Since it is asynchronous we get better result by printing it before, otherwise we always get 1 less
			server.ClientPool.Broadcast <- fmt.Sprintf("slaves %v", len(server.SlavePool.Clients)+1)
			server.SlavePool.Register <- client
			client.Read()
		} else {
			log.Println("new client")
			client := &websocket.Client{
				Conn:    conn,
				Pool:    server.ClientPool,
				Server:  server,
				IsSlave: false,
			}

			server.ClientPool.Register <- client
			client.Write(fmt.Sprintf("slaves %v", len(server.SlavePool.Clients)))
			client.Read()
		}
	})
}
