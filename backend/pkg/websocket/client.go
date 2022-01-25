package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type Client struct {
	Identifier           int
	Conn                 *websocket.Conn
	IsSlave              bool
	mu                   sync.Mutex
	CurrentSearchRequest *SearchRequest
	Server               *Server
}

func (client *Client) String() string {
	return fmt.Sprintf("[Id: %v, IsSlave: %v]", client.Identifier, client.IsSlave)
}

func (client *Client) Write(message string) {
	client.mu.Lock()
	client.Conn.WriteMessage(websocket.TextMessage, []byte(message))
	client.mu.Unlock()
}

func (client *Client) Read() {
	defer func() {
		client.Conn.Close()
	}()

	for {
		_, p, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println("client.go", client, err)
			return
		}
		message := string(p)
		messageSplit := strings.Split(message, " ")

		log.Println("client.go", "received message from user", client, ":", message)
		if client.IsSlave {
			switch command := messageSplit[0]; {
			//
			//	FOUND
			//
			case command == "found" && len(messageSplit) == 3:
				client.Server.Found(client.CurrentSearchRequest, messageSplit[2])
			}
		} else {
			switch command := messageSplit[0]; {
			//
			//	SEARCH
			//
			case command == "search" && len(messageSplit) == 2:
				re := regexp.MustCompile("^[0-9a-fA-F]{32}$")
				if re.MatchString(messageSplit[1]) {
					client.Server.AddToQueue(NewSearchRequest(messageSplit[1], client))
					client.Write(fmt.Sprintf("Added Hash %v to queue", messageSplit[1]))
				} else {
					client.Write("Wrong hash given")
				}
			//
			//	SLAVES
			//
			case command == "slaves" && len(messageSplit) == 2:
				if value, err := strconv.Atoi(messageSplit[1]); err == nil {
					if value < 1 || value > 16 {
						client.Write("Please send a value between 1 and 16")
					} else {
						if err := client.Server.Scale(value); err != nil {
							client.Write("An error occurred while trying to scale the application.")
							log.Println("client.go", client, err)
						}
					}
				} else {
					client.Write("Please give a number as second argument")
				}
			//
			//	COMMAND NOT FOUND
			//
			default:
				client.Write(fmt.Sprintf("Command \"%v\" with %v arguments not found", string(p), len(messageSplit)-1))
			}
		}
	}
}
