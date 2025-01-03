package memory

import (
	"sync"

	"github.com/Leonardo-Antonio/chemaro/dto"
)

type service struct {
	mutex   *sync.RWMutex
	storage *dto.Storage
}

func New(
	mutex *sync.RWMutex,
	storage *dto.Storage,
) *service {
	return &service{
		mutex:   mutex,
		storage: storage,
	}
}
