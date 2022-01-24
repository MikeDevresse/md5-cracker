package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"regexp"
	"strings"
	"sync"
)

type Client struct {
	ID      string
	Conn    *websocket.Conn
	Pool    *Pool
	Server  *Server
	IsSlave bool
	mu      sync.Mutex
}

// Read function that will handle all message that have been sent by the user
func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// Lock the message reading for preventing concurrency
		c.mu.Lock()

		// Split the command in an array to check part by part
		commandDetails := strings.Split(string(p), " ")
		log.Println("Received ", string(p), c.IsSlave)
		// Slave commands are different from client's one
		if c.IsSlave {
			// found hash response
			// Tell that the hash has been found and give the response
			if commandDetails[0] == "found" && len(commandDetails) == 3 {
				// Tell the other slaves to stop
				c.Server.SlavePool.Broadcast <- "stop"
				// Sends the response to the concerned client
				if c.Server.Searching != nil {
					c.Server.Searching.Client.Conn.WriteMessage(
						websocket.TextMessage,
						p,
					)
				}
				c.Server.Searching = nil
			}
		} else {
			// search hash
			// add the hash to the search waiting list
			if commandDetails[0] == "search" && len(commandDetails) == 2 {
				// Check that the given hash is md5
				re := regexp.MustCompile("^[0-9a-fA-F]{32}$")
				if re.MatchString(commandDetails[1]) {
					// Add the request to the queue
					c.Server.Queue.PushBack(&SearchRequest{
						Client:  c,
						Request: commandDetails[1],
					})
					c.Conn.WriteMessage(
						websocket.TextMessage,
						[]byte(fmt.Sprintf("searching %v", commandDetails[1])),
					)
					c.Server.ClientPool.Broadcast <- fmt.Sprintf("queue %v", c.Server.Queue.Len())
				} else {
					c.Conn.WriteMessage(
						websocket.TextMessage,
						[]byte("Wrong hash given"),
					)
				}
			} else {
				// Warn the user if no command has been found
				c.Conn.WriteMessage(
					websocket.TextMessage,
					[]byte(fmt.Sprintf("Command \"%v\" not found", string(p))),
				)
			}
		}
		c.mu.Unlock()
	}
}
