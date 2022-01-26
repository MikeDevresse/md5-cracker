package websocket

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Client struct {
	Identifier           int             // Unique Identifier, used for debug
	Conn                 *websocket.Conn // Websocket connection
	IsSlave              bool
	mu                   sync.Mutex     // Mutex that allow to lock the writing to prevent concurrency
	CurrentSearchRequest *SearchRequest // If it is a slave, the current search request that he is working on
	Server               *Server
}

// Write sends a message through the websocket connection
func (client *Client) Write(message string) {
	// We lock the mutex in order to prevent concurrency
	client.mu.Lock()
	client.Conn.WriteMessage(websocket.TextMessage, []byte(message))
	client.mu.Unlock()
}

// Read gets all the messages that the client send
func (client *Client) Read() {
	// On quit, we close the connection and remove it from the list of clients/slaves
	defer func() {
		if client.IsSlave {
			client.Server.RemoveSlave(client)
		} else {
			client.Server.RemoveClient(client)
		}
		client.Conn.Close()
	}()

	for {
		_, p, err := client.Conn.ReadMessage()
		// If an error occurred then it means that the connection is closed, we close the function
		if err != nil {
			return
		}
		message := string(p)
		messageSplit := strings.Split(message, " ")

		if client.IsSlave {
			switch command := messageSplit[0]; {
			//
			//	FOUND hash response
			//
			case command == "found" && len(messageSplit) == 3:
				client.Server.Found(client.CurrentSearchRequest, messageSplit[2])
			}
		} else {
			switch command := messageSplit[0]; {
			//
			//	SEARCH hash
			//
			case command == "search" && len(messageSplit) == 2:
				// Checks that the given hash is actually md5
				re := regexp.MustCompile("^[0-9a-fA-F]{32}$")
				if re.MatchString(messageSplit[1]) {
					// Check that the requested hash is not already present in redis
					val, err := client.Server.redis.Get(client.Server.redisContext, messageSplit[1]).Result()
					client.Write(fmt.Sprintf("Added Hash %v to queue", messageSplit[1]))
					if err == redis.Nil {
						client.Server.AddToQueue(NewSearchRequest(messageSplit[1], client))
					} else {
						searchRequest := NewSearchRequest(messageSplit[1], client)
						searchRequest.StartedAt = time.Now()
						client.Server.Found(searchRequest, val)
					}
				} else {
					client.Write("Error Please provide a valid md5 hash")
				}

			//
			//	STOP-ALL
			//
			case command == "stop-all" && len(messageSplit) == 1:
				client.Server.StopAll()

			//
			//	AUTO-SCALE true|false
			//
			case command == "auto-scale" && len(messageSplit) == 2:
				client.Server.SetAutoScale(messageSplit[1] == "true")

			//
			//  MAX-SEARCH maxSearchParameter
			//
			case command == "max-search" && len(messageSplit) == 2:
				// We want the max search to be at the end of the alphabet, so we ask that it is only 9 from 2 to 8 times
				re := regexp.MustCompile("^[9]{2,8}$")
				if re.MatchString(messageSplit[1]) {
					client.Server.SetMaxSearch(messageSplit[1])
				} else {
					client.Write("Error max-search argument must follow the regex: ^[9]{2,8}$")
				}

			//
			//  MAX-SLAVES-PER-REQUEST number
			//
			case command == "max-slaves-per-request" && len(messageSplit) == 2:
				// Check that the given parameter is an int
				if value, err := strconv.Atoi(messageSplit[1]); err == nil {
					if value < 1 {
						client.Write("Error max-slaves-per-request must be greater than 0")
					} else {
						client.Server.SetMaxSlavesPerRequest(value)
					}
				} else {
					client.Write("Error max-slaves-per-request expects a number as second parameter")
				}

			//
			//	SLAVES number
			//
			case command == "slaves" && len(messageSplit) == 2:
				// Check that the given parameter is an int
				if value, err := strconv.Atoi(messageSplit[1]); err == nil {
					if value < 1 || value > 16 {
						client.Write("Error slaves must be between 1 and 16")
					} else {
						client.Write("Scaling")
						// If an error occurred while scaling the slaves
						if err := client.Server.Scale(value); err != nil {
							client.Write("Error An error occurred while trying to scale the application.")
						}
					}
				} else {
					client.Write("Error slaves expects a number as second parameter")
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
