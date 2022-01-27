package main

import (
	"fmt"
	"github.com/MikeDevresse/md5-cracker/backend/pkg/websocket"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"os"
	"strconv"
)

// main init the websocket and starts the webserver
func main() {
	port := getEnv("SERVER_PORT", "80")
	initWebsocket()
	log.Println("main.go", "Running webserver on port", port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		log.Fatal("An error occurred while trying to start the webserver: \"", err, "\"")
	}
}

// initWebsocket initialize the websocket route and the server
func initWebsocket() {
	db, err := strconv.Atoi(getEnv("REDIS_DATABASE", "0"))
	if err != nil {
		log.Fatal("REDIS_DATABASE environment variable must be an integer")
	}
	server := websocket.NewServer(redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_HOST", "redis:6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		Username: getEnv("REDIS_USERNAME", ""),
		DB:       db,
	}))
	go server.Start()

	clientCount := 0
	http.HandleFunc(getEnv("SERVER_PATH", "/ws"), func(w http.ResponseWriter, r *http.Request) {
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

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
