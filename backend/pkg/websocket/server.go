package websocket

import (
	"container/list"
	"context"
	"fmt"
	"github.com/MikeDevresse/md5-cracker/backend/pkg/service"
	"github.com/go-redis/redis/v8"
	"os/exec"
	"regexp"
	"time"
)

type Server struct {
	Slaves              map[*Client]bool        // List of Clients that are slaves
	AvailableSlaves     map[*Client]bool        // List of slaves that are available for work
	Clients             map[*Client]bool        // List of connected Clients
	Queue               *list.List              // Queue containing the SearchRequest that are waiting to be resolved
	Searching           map[*SearchRequest]bool // List of SearchRequest that are being searched
	MaxSearch           string                  // The limit at which the message would be, can only be 9 from 2 to 8 times
	MaxSlavesPerRequest int                     // Maximum of slaves that will be dedicated to a SearchRequest
	SlavesCount         int                     // Number of slaves that we want (not equals to len(Slaves) if we are currently scaling)
	AutoScale           bool                    // Tells whether the app should scale the slave automatically depending on the queue
	redis               *redis.Client           // Redis connection
	redisContext        context.Context         // context that is used for redis
}

// NewServer initialized a Server object with default values and parameters
func NewServer(client *redis.Client) *Server {
	return &Server{
		Slaves:              make(map[*Client]bool),
		AvailableSlaves:     make(map[*Client]bool),
		Clients:             make(map[*Client]bool),
		Queue:               list.New(),
		Searching:           make(map[*SearchRequest]bool),
		MaxSearch:           "9999",
		MaxSlavesPerRequest: 4,
		SlavesCount:         0,
		AutoScale:           false,
		redis:               client,
		redisContext:        context.Background(),
	}
}

// AddSlave adds a slave to the Slaves list and the AvailableSlaves list, and calls BroadcastSlaveStatus
func (server *Server) AddSlave(slave *Client) {
	server.Slaves[slave] = true
	server.AvailableSlaves[slave] = true
	server.BroadcastSlaveStatus()
}

// RemoveSlave removes a slave from the Slaves and the AvailableSlaves list, and calls BroadcastSlaveStatus
func (server *Server) RemoveSlave(slave *Client) {
	delete(server.Slaves, slave)
	delete(server.AvailableSlaves, slave)
	server.BroadcastSlaveStatus()
}

// AddClient adds a client to the Clients list
func (server *Server) AddClient(client *Client) {
	server.Clients[client] = true
}

// RemoveClient removes a client from the Clients list
func (server *Server) RemoveClient(client *Client) {
	delete(server.Clients, client)
}

// AddToQueue adds a search request to the queue and calls BroadcastQueueStatus
func (server *Server) AddToQueue(request *SearchRequest) {
	server.Queue.PushBack(request)
	server.BroadcastQueueStatus()
}

// BroadcastQueueStatus sends to all Clients the current status of the queue in the format:
// queue numberOfRequestInQueue numberOfRequestBeingHandled
func (server *Server) BroadcastQueueStatus() {
	message := fmt.Sprintf("queue %v %v", server.Queue.Len(), len(server.Searching))
	for client := range server.Clients {
		client.Write(message)
	}
}

// BroadcastSlaveStatus sends to all Clients the current status of slaves in the format:
// slaves SlavesCount numberOfSlavesNotWorking numberOfSlavesWorking
func (server *Server) BroadcastSlaveStatus() {
	message := fmt.Sprintf("slaves %v %v %v", server.SlavesCount, len(server.AvailableSlaves), len(server.Slaves)-len(server.AvailableSlaves))
	for client := range server.Clients {
		client.Write(message)
	}
}

// BroadcastConfiguration sends to all Clients the current configuration in the format:
// max-search MaxSearch
// slaves SlavesCount numberOfSlavesNotWorking numberOfSlavesWorking
// max-slaves-per-request MaxSlavesPerRequest
// auto-scale true|false
func (server *Server) BroadcastConfiguration() {
	for client := range server.Clients {
		client.Write(fmt.Sprintf("max-search %v", server.MaxSearch))
		client.Write(fmt.Sprintf("slaves %v %v %v", len(server.Slaves), len(server.AvailableSlaves), len(server.Slaves)-len(server.AvailableSlaves)))
		client.Write(fmt.Sprintf("max-slaves-per-request %v", server.MaxSlavesPerRequest))
		client.Write(fmt.Sprintf("auto-scale %v", server.AutoScale))
	}
}

// PrintConfiguration sends the current configuration to a client in the format
// max-search MaxSearch
// slaves SlavesCount numberOfSlavesNotWorking numberOfSlavesWorking
// max-slaves-per-request MaxSlavesPerRequest
// auto-scale true|false
func (server *Server) PrintConfiguration(client *Client) {
	client.Write(fmt.Sprintf("max-search %v", server.MaxSearch))
	client.Write(fmt.Sprintf("slaves %v %v %v", len(server.Slaves), len(server.AvailableSlaves), len(server.Slaves)-len(server.AvailableSlaves)))
	client.Write(fmt.Sprintf("max-slaves-per-request %v", server.MaxSlavesPerRequest))
	client.Write(fmt.Sprintf("auto-scale %v", server.AutoScale))
}

// SetMaxSearch set the MaxSearch parameter, it defines to what limit the slaves must search the word, it can only be
// 9 from 2 to 8 times
func (server *Server) SetMaxSearch(maxSearch string) {
	re := regexp.MustCompile("^[9]{2,8}$")
	if !re.MatchString(maxSearch) || server.MaxSearch == maxSearch {
		return
	}
	server.MaxSearch = maxSearch
	server.BroadcastConfiguration()
}

// SetMaxSlavesPerRequest set the maximum number of slaves that we assign to a SearchRequest
func (server *Server) SetMaxSlavesPerRequest(maxSlavesPerRequest int) {
	if maxSlavesPerRequest < 1 || server.MaxSlavesPerRequest == maxSlavesPerRequest {
		return
	}
	server.MaxSlavesPerRequest = maxSlavesPerRequest
	server.BroadcastConfiguration()
}

// SetAutoScale set if the server should scale automatically the slaves or not
func (server *Server) SetAutoScale(isAutoScale bool) {
	if server.AutoScale == isAutoScale {
		return
	}
	server.AutoScale = isAutoScale
	server.BroadcastConfiguration()
}

// Found called when a slave has found the result of a SearchRequest, it then sends the result to the client and tells
// the other Slaves that were working on the SearchRequest to stop
func (server *Server) Found(request *SearchRequest, result string) {
	server.redis.Set(server.redisContext, request.Hash, result, 0)
	request.Result = result
	request.EndedAt = time.Now()
	request.Client.Write(fmt.Sprintf("found %v", request.FormatResponse()))
	delete(server.Searching, request)
	for slave := range request.Slaves {
		slave.Write("stop")
		server.AvailableSlaves[slave] = true
	}
	server.BroadcastQueueStatus()
	server.BroadcastSlaveStatus()
}

// StopAll stop all the SearchRequest that are running, clear the queue, and tells the Slaves to stop
func (server *Server) StopAll() {
	for slave := range server.Slaves {
		slave.Write("stop")
		server.AvailableSlaves[slave] = true
	}
	server.Queue = list.New()
	server.Searching = make(map[*SearchRequest]bool)
	server.BroadcastQueueStatus()
	server.BroadcastSlaveStatus()
}

// Scale choose the number of Slaves that we want in our application, must be between 0 and 16 for performance reason
func (server *Server) Scale(number int) error {
	if number < 0 {
		number = 0
	}
	if number > 16 {
		number = 16
	}
	server.SlavesCount = number
	// If we do not have enough Slaves then we call the docker-compose up -d --scale command that allow us to add
	// more slaves
	if len(server.Slaves) < server.SlavesCount {
		cmd := exec.Command("docker-compose", "up", "-d", "--no-recreate", "--scale", fmt.Sprintf("slave=%v", number))
		return cmd.Run()
	}

	return nil
}

// Start Loop that will handle autoscaling, soft downscale, and queue handling
func (server *Server) Start() {
	for {
		// If we are on autoScale mode then we always want MaxSlavesPerRequest  slaves available for work with a maximum
		//of 16 slaves
		if server.AutoScale {
			server.Scale(service.Min(16, server.MaxSlavesPerRequest*(server.Queue.Len()+1)))
		}
		// If we have too many Slaves compared to the amount asked, then we wait for slaves to be available in order to
		// tell them to exit
		if len(server.Slaves) > server.SlavesCount && len(server.AvailableSlaves) != 0 {
			toDelete := len(server.Slaves) - server.SlavesCount
			for slave := range server.AvailableSlaves {
				slave.Write("exit")
				delete(server.Slaves, slave)
				delete(server.AvailableSlaves, slave)
				toDelete--
				if toDelete <= 0 {
					break
				}
			}
		}
		// If we have elements in the Queue, and we have slaves not working then we handle an element of the queue
		if server.Queue.Len() != 0 && len(server.AvailableSlaves) != 0 {
			// Dequeue the element
			elem := server.Queue.Front()
			server.Queue.Remove(elem)
			// Cast the element to SearchRequest
			searchRequest := elem.Value.(*SearchRequest)
			searchRequest.StartedAt = time.Now()
			// Tells that we are currently searching for this
			server.Searching[searchRequest] = true

			// Get the amount of slaves that will be working for this request
			maxSlavePerRequest := server.MaxSlavesPerRequest
			slaveCount := service.Min(maxSlavePerRequest, len(server.AvailableSlaves))

			// Number of possibility that will a slave calculate
			division := float64(service.Convert62to10(server.MaxSearch)) / float64(slaveCount)
			i := 0
			// Divide the task for each slave connected
			toRemove := make([]*Client, 0)
			for slave := range server.AvailableSlaves {
				slave.Write(fmt.Sprintf(
					"search %v %v %v",
					searchRequest.Hash,
					service.Convert10to62(int(division*float64(i))),
					service.Convert10to62(int(division*float64(i+1))),
				))
				toRemove = append(toRemove, slave)
				i = i + 1
				if i >= slaveCount {
					break
				}
			}

			// Remove the slave from being available
			for _, slave := range toRemove {
				slave.CurrentSearchRequest = searchRequest
				searchRequest.Slaves[slave] = true
				delete(server.AvailableSlaves, slave)
			}

			server.BroadcastQueueStatus()
			server.BroadcastSlaveStatus()
		}
	}
}
