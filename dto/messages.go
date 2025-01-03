package dto

type Storage struct {
	Pool map[string][]Message `json:"pool"`
}

type Message struct {
	Id        string `json:"id"`
	UserId    string `json:"userId"`
	Message   string `json:"message"`
	Type      string `json:"type"`
	CreatedAt uint64 `json:"createdAt"`
}
