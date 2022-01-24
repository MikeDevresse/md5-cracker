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

		commandDetails := strings.Split(string(p), " ")
		log.Println("Received ", string(p), c.IsSlave)
		if c.IsSlave {
			if commandDetails[0] == "found" && len(commandDetails) == 3 {
				c.Server.SlavePool.Broadcast <- "stop"
				if c.Server.Searching != nil {
					c.Server.Searching.Client.mu.Lock()
					c.Server.Searching.Client.Conn.WriteMessage(
						websocket.TextMessage,
						p,
					)
					c.Server.Searching.Client.mu.Unlock()
				}
				c.Server.Searching = nil
			}
		} else {
			if commandDetails[0] == "search" && len(commandDetails) == 2 {
				re := regexp.MustCompile("^[0-9a-fA-F]{32}$")
				c.mu.Lock()
				if re.MatchString(commandDetails[1]) {
					c.Server.Queue.PushBack(&SearchRequest{
						Client:  c,
						Request: commandDetails[1],
					})
					c.Conn.WriteMessage(
						websocket.TextMessage,
						[]byte(fmt.Sprintf("searching %v", commandDetails[1])),
					)
				} else {
					c.Conn.WriteMessage(
						websocket.TextMessage,
						[]byte("Wrong hash given"),
					)
				}
				c.mu.Unlock()

			} else {
				c.mu.Lock()
				c.Conn.WriteMessage(
					websocket.TextMessage,
					[]byte(fmt.Sprintf("Command \"%v\" not found", string(p))),
				)
				c.mu.Unlock()
			}
		}
	}
}
