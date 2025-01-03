package memory

import "github.com/Leonardo-Antonio/chemaro/dto"

func (_s *service) Get(groupId string) []dto.Message {
	_s.mutex.RLock()
	found := _s.storage.Pool[groupId]
	_s.mutex.RUnlock()
	return found
}

func (_s *service) GetAll() map[string][]dto.Message {
	_s.mutex.RLock()
	if _s.storage.Pool == nil {
		return map[string][]dto.Message{}
	}
	found := _s.storage.Pool
	_s.mutex.RUnlock()
	return found
}
