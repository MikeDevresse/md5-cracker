package websocket

import (
	"container/list"
	"fmt"
	"github.com/MikeDevresse/md5-cracker/pkg/service"
	"github.com/gorilla/websocket"
	"log"
)

type Server struct {
	Queue      *list.List
	Searching  *SearchRequest
	ClientPool *Pool
	SlavePool  *Pool
}

func NewServer() *Server {
	return &Server{
		Queue:      list.New(),
		ClientPool: NewPool(),
		SlavePool:  NewPool(),
	}
}

func (server *Server) Start() {
	go server.SlavePool.Start()
	go server.ClientPool.Start()
	for {
		if server.Queue.Len() != 0 && len(server.SlavePool.Clients) > 0 && server.Searching == nil {
			elem := server.Queue.Front()
			toSearch := elem.Value.(*SearchRequest)
			server.Queue.Remove(elem)
			server.Searching = toSearch
			division := service.Convert62to10("9999") / len(server.SlavePool.Clients)
			i := 0
			for slave := range server.SlavePool.Clients {
				req := fmt.Sprintf(
					"search %v %v %v",
					toSearch.Request,
					service.Convert10to62(division*i),
					service.Convert10to62(division*(i+1)),
				)
				log.Println(req)
				slave.mu.Lock()
				slave.Conn.WriteMessage(
					websocket.TextMessage,
					[]byte(req),
				)
				slave.mu.Unlock()
				i = i + 1
			}
		}
	}
}
