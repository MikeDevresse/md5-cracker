package websocket

import (
	"container/list"
	"fmt"
	"github.com/MikeDevresse/md5-cracker/pkg/service"
	"log"
)

type Server struct {
	Queue      *list.List     // List of SearchRequest to search
	Searching  *SearchRequest // Current SearchRequest
	ClientPool *Pool          // Pool containing only normal clients
	SlavePool  *Pool          // Pool containing only slaves
}

func NewServer() *Server {
	return &Server{
		Queue:      list.New(),
		ClientPool: NewPool(),
		SlavePool:  NewPool(),
	}
}

// Start function that tells the slaves to search for a hash if one is present in the queue
func (server *Server) Start() {
	go server.SlavePool.Start()
	go server.ClientPool.Start()
	for {
		// If the queue is not empty, that we have slaves connected and we are not currently searching then send
		// the commands to the slaves
		if server.Queue.Len() != 0 && len(server.SlavePool.Clients) > 0 && server.Searching == nil {
			// Dequeue
			elem := server.Queue.Front()
			server.Queue.Remove(elem)
			server.Searching = elem.Value.(*SearchRequest)
			// Number of possibility that will a slave calculate
			division := float64(service.Convert62to10("9999")) / float64(len(server.SlavePool.Clients))
			i := 0
			// Divide the task for each slave connected
			for slave := range server.SlavePool.Clients {
				req := fmt.Sprintf(
					"search %v %v %v",
					server.Searching.Request,
					service.Convert10to62(int(division)*i),
					service.Convert10to62(int(division)*(i+1)),
				)
				log.Println(req)
				slave.Write(req)
				i = i + 1
			}
		}
	}
}
