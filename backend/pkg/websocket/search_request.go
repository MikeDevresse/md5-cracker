package websocket

import "fmt"

type SearchRequest struct {
	Hash   string
	Slaves map[*Client]bool
	Client *Client
}

func (searchRequest *SearchRequest) String() string {
	return fmt.Sprintf("[Client: %v, Hash: %v]", searchRequest.Client, searchRequest.Hash)
}

func NewSearchRequest(hash string, client *Client) *SearchRequest {
	return &SearchRequest{
		Hash:   hash,
		Slaves: make(map[*Client]bool),
		Client: client,
	}
}
