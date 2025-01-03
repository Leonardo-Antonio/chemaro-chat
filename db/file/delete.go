package memory

import (
	"log"
	"os"

	"github.com/Leonardo-Antonio/chemaro/dto"
)

func (_s *service) Delete(groupId string) {
	_s.mutex.Lock()
	delete(_s.storage.Pool, groupId)
	_s.mutex.Unlock()
}

func (_s *service) DeleteAll() {
	_s.mutex.Lock()
	_s.storage.Pool = make(map[string][]dto.Message)
	os.WriteFile("storage.json", []byte("{\"pool\":{}}"), 0644)
	log.Println("Deleted all messages")
	_s.mutex.Unlock()
}
