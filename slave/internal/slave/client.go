package slave

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/MikeDevresse/md5-cracker/slave/pkg/bytes_incrementor"
	"github.com/gorilla/websocket"
	"log"
	"strings"
	"sync"
)

type client struct {
	conn  *websocket.Conn
	mu    sync.Mutex
	done  chan bool
	Write chan string
}

func NewClient(conn *websocket.Conn) *client {
	return &client{
		conn:  conn,
		done:  make(chan bool),
		Write: make(chan string),
	}
}

func (client *client) Start() {
	go client.Read()
	// Tells the master that we are a slave
	client.conn.WriteMessage(websocket.TextMessage, []byte("slave"))
	for {
		select {
		case <-client.done:
			return
		case message := <-client.Write:
			log.Println("Writing ", message)
			client.mu.Lock()
			client.conn.WriteMessage(websocket.TextMessage, []byte(message))
			client.mu.Unlock()
			log.Println("Done writing ", message)
		}
	}
}

func (client *client) Read() {
	stop := make(chan bool)
	running := false
	for {
		_, p, _ := client.conn.ReadMessage()
		message := string(p)
		log.Println("Received", message)
		messageSplit := strings.Split(message, " ")

		switch command := messageSplit[0]; {
		case command == "search" && len(messageSplit) == 4:
			running = true
			go func() {
				hash := messageSplit[1]
				encodedHash, _ := hex.DecodeString(hash)
				start := []byte(messageSplit[2])
				end := []byte(messageSplit[3])
				for running {
					select {
					case <-stop:
						running = false
						return
					default:
						md5sum := md5.Sum(start)
						if bytes.Equal(md5sum[:], encodedHash) {
							client.Write <- fmt.Sprintf("found %v %v", hash, string(start))
							running = false
						} else {
							bytes_incrementor.Increment(&start, len(start)-1)
							if bytes.Equal(start, end) {
								<-stop
							}
						}
					}
				}
			}()
		case command == "stop":
			if running == true {
				stop <- true
			}
		case command == "exit":
			client.done <- true
		}
	}
}
