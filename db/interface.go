package db

import (
	"sync"

	filedb "github.com/Leonardo-Antonio/chemaro/db/file"
	"github.com/Leonardo-Antonio/chemaro/db/memory"
	"github.com/Leonardo-Antonio/chemaro/dto"
)

type (
	IService interface {
		DeleteAll()
		Delete(groupId string)
		Get(groupId string) []dto.Message
		GetAll() map[string][]dto.Message
		Set(groupId string, message dto.Message)
	}

	driver string
)

const (
	DB_MEMORY = "memory"
	DB_REDIS  = "redis"
	DB_FILE   = "file"
)

var DB IService

func New(driver driver, mx *sync.RWMutex, storage *dto.Storage) IService {
	switch driver {
	case DB_MEMORY:
		return memory.New(mx, storage)
	case DB_FILE:
		return filedb.New(mx, storage)
	default:
		return memory.New(mx, storage)
	}
}
