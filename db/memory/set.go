package memory

func Set(groupId string, message Message) {
	mutex.Lock()
	if storage.Pool == nil {
		storage.Pool = make(map[string][]Message)
	}
	storage.Pool[groupId] = append(storage.Pool[groupId], message)
	mutex.Unlock()
}
