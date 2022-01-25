package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"os/exec"
	"regexp"
	"strconv"
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

func (c *Client) Write(msg string) {
	c.mu.Lock()
	c.Conn.WriteMessage(
		websocket.TextMessage,
		[]byte(msg),
	)
	c.mu.Unlock()
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

		// Split the command in an array to check part by part
		commandDetails := strings.Split(string(p), " ")
		log.Println("Received ", string(p), c.IsSlave)
		// Slave commands are different from client's one
		if c.IsSlave {
			// found hash response
			// Tell that the hash has been found and give the response
			if commandDetails[0] == "found" && len(commandDetails) == 3 {
				// Tell the other slaves to stop (can't use broadcast otherwise sometimes the new search comes before than
				// the stop one, and it gets stuck)
				for slave := range c.Server.SlavePool.Clients {
					slave.Write("stop")
				}
				// Sends the response to the concerned client
				if c.Server.Searching != nil {
					c.Server.Searching.Client.Write(string(p))
				}
				c.Server.ClientPool.Broadcast <- fmt.Sprintf("queue %v", c.Server.Queue.Len())
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
					// TODO: Fix bug that tells the queue is at 2 but it is at 1
					c.Server.ClientPool.Broadcast <- fmt.Sprintf("queue %v", c.Server.Queue.Len()+1)
					c.Write(fmt.Sprintf("searching %v", commandDetails[1]))
				} else {
					c.Write("Wrong hash given")
				}
			} else if commandDetails[0] == "slaves" && len(commandDetails) == 2 {
				if value, err := strconv.Atoi(commandDetails[1]); err == nil {
					if value < 1 || value > 16 {
						c.Write("Please send a value between 1 and 16")
					} else {
						cmd := exec.Command(
							"docker-compose",
							"up",
							"-d",
							"--no-recreate",
							"--scale",
							fmt.Sprintf("slave=%v", value),
						)

						if err := cmd.Run(); err != nil {
							c.Write("An error occured while trying to scale the application.")
							log.Println(err)
						}
					}
				} else {
					c.Write("Please give a number as second argument")
				}
			} else {
				// Warn the user if no command has been found
				c.Write(fmt.Sprintf("Command \"%v\" not found", string(p)))
			}
		}
	}
}
