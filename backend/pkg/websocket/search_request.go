package websocket

import (
	"fmt"
	"strconv"
	"time"
)

type SearchRequest struct {
	Hash        string           // Hash that is being searched
	Slaves      map[*Client]bool // Slaves that are working one resolving the hash
	Client      *Client          // Client that made the request
	Result      string           // Decoded hash
	RequestedAt time.Time        // Time when the request has been received
	StartedAt   time.Time        // Time when the slaves started to work on this request
	EndedAt     time.Time        // Time when the result has been found
}

// FormatResponse format the object in order to be returned in a response to a Client
func (searchRequest *SearchRequest) FormatResponse() string {
	return fmt.Sprintf(
		"%v %v %v %v %v",
		searchRequest.Hash,
		searchRequest.Result,
		strconv.FormatInt(searchRequest.RequestedAt.UnixMilli(), 10),
		strconv.FormatInt(searchRequest.StartedAt.UnixMilli(), 10),
		strconv.FormatInt(searchRequest.EndedAt.UnixMilli(), 10),
	)
}

// NewSearchRequest initialize a SearchRequest object with default values and given parameters
func NewSearchRequest(hash string, client *Client) *SearchRequest {
	return &SearchRequest{
		Hash:        hash,
		Slaves:      make(map[*Client]bool),
		Client:      client,
		RequestedAt: time.Now(),
	}
}
