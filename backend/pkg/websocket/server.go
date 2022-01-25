package websocket

import (
	"container/list"
	"fmt"
	"github.com/MikeDevresse/md5-cracker/pkg/service"
	"log"
	"os/exec"
)

type Server struct {
	Slaves          map[*Client]bool
	AvailableSlaves map[*Client]bool
	Clients         map[*Client]bool
	Queue           *list.List
	Searching       map[*SearchRequest]bool
}

func NewServer() *Server {
	return &Server{
		Slaves:          make(map[*Client]bool),
		AvailableSlaves: make(map[*Client]bool),
		Clients:         make(map[*Client]bool),
		Queue:           list.New(),
		Searching:       make(map[*SearchRequest]bool),
	}
}

func (server *Server) AddSlave(slave *Client) {
	server.Slaves[slave] = true
	server.AvailableSlaves[slave] = true
	server.BroadcastSlaveStatus()
}

func (server *Server) RemoveSlave(slave *Client) {
	delete(server.Slaves, slave)
	delete(server.AvailableSlaves, slave)
	server.BroadcastSlaveStatus()
}

func (server *Server) AddClient(client *Client) {
	server.Clients[client] = true
}

func (server *Server) RemoveClient(client *Client) {
	delete(server.Clients, client)
}

func (server *Server) AddToQueue(request *SearchRequest) {
	server.Queue.PushBack(request)
	server.BroadcastQueueStatus()
	log.Println("server.go", "New Request:", request)
}

func (server *Server) BroadcastQueueStatus() {
	message := fmt.Sprintf("queue %v %v", server.Queue.Len(), len(server.Searching))
	for client := range server.Clients {
		client.Write(message)
	}
}

func (server *Server) BroadcastSlaveStatus() {
	message := fmt.Sprintf("slaves %v %v", len(server.AvailableSlaves), len(server.Slaves)-len(server.AvailableSlaves))
	for client := range server.Clients {
		client.Write(message)
	}
}

func (server *Server) Found(request *SearchRequest, result string) {
	request.Client.Write(fmt.Sprintf("found %v %v", request.Hash, result))
	delete(server.Searching, request)
	for slave := range request.Slaves {
		slave.Write("stop")
		server.AvailableSlaves[slave] = true
	}
	server.BroadcastQueueStatus()
	server.BroadcastSlaveStatus()
}

func (server *Server) Scale(number int) error {
	if number < 0 {
		number = 0
	}
	if number > 16 {
		number = 16
	}
	cmd := exec.Command("docker-compose", "up", "-d", "--no-recreate", "--scale", fmt.Sprintf("slave=%v", number))
	err := cmd.Run()

	return err
}

func (server *Server) Start() {
	for {
		if server.Queue.Len() != 0 && len(server.AvailableSlaves) != 0 {
			// Dequeue the element
			elem := server.Queue.Front()
			server.Queue.Remove(elem)
			// Cast the element to SearchRequest
			searchRequest := elem.Value.(*SearchRequest)
			// Tells that we are currently searching for this
			server.Searching[searchRequest] = true

			// Get the amount of slaves that will be working for this request
			maxSlavePerRequest := 4
			slaveCount := service.Min(maxSlavePerRequest, len(server.AvailableSlaves))

			// Number of possibility that will a slave calculate
			division := float64(service.Convert62to10("9999")) / float64(slaveCount)
			i := 0
			// Divide the task for each slave connected
			toRemove := make([]*Client, 0)
			for slave := range server.AvailableSlaves {
				req := fmt.Sprintf(
					"search %v %v %v",
					searchRequest.Hash,
					service.Convert10to62(int(division*float64(i))),
					service.Convert10to62(int(division*float64(i+1))),
				)
				log.Println("server.go", "Sending", req, "to", slave)
				slave.Write(req)
				toRemove = append(toRemove, slave)
				i = i + 1
				if i >= slaveCount {
					break
				}
			}

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
