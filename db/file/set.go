package memory

import (
	"encoding/json"
	"os"

	"github.com/Leonardo-Antonio/chemaro/dto"
)

func (_s *service) Set(groupId string, message dto.Message) {
	_s.mutex.Lock()
	buff, err := os.ReadFile("storage.json")
	if err != nil {
		panic(err)
	}

	var storage dto.Storage
	err = json.Unmarshal(buff, &storage)
	if err != nil {
		panic(err)
	}

	storage.Pool[groupId] = append(storage.Pool[groupId], message)

	buff, err = json.Marshal(storage)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("storage.json", buff, 0644)
	if err != nil {
		panic(err)
	}
	_s.mutex.Unlock()
}
