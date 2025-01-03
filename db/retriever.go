package db

import (
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"github.com/Leonardo-Antonio/chemaro/dto"
)

var (
	mutex   = new(sync.RWMutex)
	storage = &dto.Storage{
		Pool: make(map[string][]dto.Message),
	}
)

func RunRetriever() {
	go func(mx *sync.RWMutex, s *dto.Storage) {
		for {
			daemon(mx, s)
		}
	}(mutex, storage)
}

func daemon(mx *sync.RWMutex, storage *dto.Storage) {
	TTL_DB_SECONS_STRING := os.Getenv("TTL_DB_SECONS")
	TTL_DB_SECONS, err := strconv.Atoi(TTL_DB_SECONS_STRING)
	if err != nil {
		TTL_DB_SECONS = 60
	}

	ticker := time.NewTicker(time.Duration(TTL_DB_SECONS) * time.Second)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recover from panic\n", string(debug.Stack()))
		}
		ticker.Stop()
	}()

	DB = New(DB_MEMORY, mx, storage)

	for {
		select {
		case <-ticker.C:
			DB.DeleteAll()
		default:
			continue
		}
	}
}
