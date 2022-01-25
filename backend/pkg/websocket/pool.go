package websocket

import "fmt"

type Pool struct {
	Clients    map[*Client]bool // Client list
	Register   chan *Client     // Adds the given Client to the Pool
	Unregister chan *Client     // Delete a Client from the Pool
	Broadcast  chan string      // Sends a message to all Client in the Pool
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan string),
	}
}

// Start Function that handle given command from a Pool
func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			if client.IsSlave {
				client.Server.ClientPool.Broadcast <- fmt.Sprintf("slaves %v", len(pool.Clients))
			}
		case message := <-pool.Broadcast:
			for client := range pool.Clients {
				client.Write(message)
			}
		}
	}
}
