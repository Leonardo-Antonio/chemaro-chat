package memory

import "github.com/Leonardo-Antonio/chemaro/dto"

func (_s *service) Set(groupId string, message dto.Message) {
	_s.mutex.Lock()
	if _s.storage.Pool == nil {
		_s.storage.Pool = make(map[string][]dto.Message)
	}
	_s.storage.Pool[groupId] = append(_s.storage.Pool[groupId], message)
	_s.mutex.Unlock()
}
