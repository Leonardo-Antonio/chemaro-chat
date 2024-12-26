package memory

import "sync"

type Message struct {
	Id        string `json:"id"`
	UserId    string `json:"userId"`
	Message   string `json:"message"`
	Type      string `json:"type"`
	CreatedAt uint64 `json:"createdAt"`
}

var (
	storage = make(map[string][]Message)
	mutex   = &sync.RWMutex{}
)
