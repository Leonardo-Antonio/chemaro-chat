package memory

import (
	"encoding/json"
	"os"

	"github.com/Leonardo-Antonio/chemaro/dto"
)

func (_s *service) Get(groupId string) []dto.Message {
	_s.mutex.RLock()
	buff, err := os.ReadFile("storage.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(buff, &_s.storage)
	if err != nil {
		panic(err)
	}

	found := _s.storage.Pool[groupId]
	_s.mutex.RUnlock()
	return found
}

func (_s *service) GetAll() map[string][]dto.Message {
	_s.mutex.RLock()
	buff, err := os.ReadFile("storage.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(buff, &_s.storage)
	if err != nil {
		panic(err)
	}

	found := _s.storage.Pool
	_s.mutex.RUnlock()
	return found
}
