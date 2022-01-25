package websocket

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
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
					val, err := client.Server.redis.Get(client.Server.redisContext, messageSplit[1]).Result()
					if err == redis.Nil {
						client.Server.AddToQueue(NewSearchRequest(messageSplit[1], client))
						client.Write(fmt.Sprintf("Added Hash %v to queue", messageSplit[1]))
					} else {
						searchRequest := NewSearchRequest(messageSplit[1], client)
						searchRequest.StartedAt = time.Now()
						client.Server.Found(searchRequest, val)
					}
				} else {
					client.Write("Wrong hash given")
				}
			//
			//	STOP-ALL
			//
			case command == "stop-all" && len(messageSplit) == 1:
				client.Server.StopAll()
			//
			//	AUTO-SCALE
			//
			case command == "auto-scale" && len(messageSplit) == 2:
				client.Server.SetAutoScale(messageSplit[1] == "true")
				client.Server.PrintConfiguration(client)
			//
			//  MAX-SEARCH
			//
			case command == "max-search" && len(messageSplit) == 2:
				re := regexp.MustCompile("^[9]{2,8}$")
				if re.MatchString(messageSplit[1]) {
					client.Server.SetMaxSearch(messageSplit[1])
					client.Server.PrintConfiguration(client)
				} else {
					client.Write("Wrong max-search argument, it must follow the regex: ^[9]{2,8}$")
				}
			//
			//  MAX-SLAVES-PER-REQUEST
			//
			case command == "max-slaves-per-request" && len(messageSplit) == 2:
				if value, err := strconv.Atoi(messageSplit[1]); err == nil {
					if value < 1 {
						client.Write("Please greater than 0")
					} else {
						client.Server.SetMaxSlavesPerRequest(value)
						client.Server.PrintConfiguration(client)
					}
				} else {
					client.Write("Please give a number as second argument")
				}
			//
			//	SLAVES
			//
			case command == "slaves" && len(messageSplit) == 2:
				if value, err := strconv.Atoi(messageSplit[1]); err == nil {
					if value < 1 || value > 16 {
						client.Write("Please send a value between 1 and 16")
					} else {
						client.Write("Scaling ...")
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
