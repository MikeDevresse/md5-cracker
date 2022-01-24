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

type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		message := Message{Type: messageType, Body: string(p)}
		commandDetails := strings.Split(message.Body, " ")
		log.Println("Received ", message.Body, c.IsSlave)
		if c.IsSlave {
			if commandDetails[0] == "found" && len(commandDetails) == 3 {
				for slave := range c.Server.SlavePool.Clients {
					slave.mu.Lock()
					slave.Conn.WriteMessage(
						websocket.TextMessage,
						[]byte("stop"),
					)
					slave.mu.Unlock()
				}
				if c.Server.Searching != nil {
					c.Server.Searching.Client.mu.Lock()
					c.Server.Searching.Client.Conn.WriteMessage(
						websocket.TextMessage,
						[]byte(message.Body),
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

			}
		}
		fmt.Printf("Message Received: %+v\n", message)
	}
}
