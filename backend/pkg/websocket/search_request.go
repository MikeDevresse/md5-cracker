package websocket

import (
	"fmt"
	"strings"
	"time"
)

type SearchRequest struct {
	Hash        string
	Slaves      map[*Client]bool
	Client      *Client
	Result      string
	RequestedAt time.Time
	StartedAt   time.Time
	EndedAt     time.Time
}

func (searchRequest *SearchRequest) String() string {
	return fmt.Sprintf("[Client: %v, Hash: %v]", searchRequest.Client, searchRequest.Hash)
}

func (searchRequest *SearchRequest) FormatResponse() string {
	return fmt.Sprintf(
		"%v %v %v %v %v",
		searchRequest.Hash,
		searchRequest.Result,
		strings.ReplaceAll(searchRequest.RequestedAt.String(), " ", "_"),
		strings.ReplaceAll(searchRequest.StartedAt.String(), " ", "_"),
		strings.ReplaceAll(searchRequest.EndedAt.String(), " ", "_"),
	)
}

func NewSearchRequest(hash string, client *Client) *SearchRequest {
	return &SearchRequest{
		Hash:        hash,
		Slaves:      make(map[*Client]bool),
		Client:      client,
		RequestedAt: time.Now(),
	}
}
