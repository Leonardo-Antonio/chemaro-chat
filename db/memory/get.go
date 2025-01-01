package memory

func Get(groupId string) []Message {
	mutex.RLock()
	found := storage.Pool[groupId]
	mutex.RUnlock()
	return found
}

func GetAll() map[string][]Message {
	mutex.RLock()
	if storage.Pool == nil {
		return map[string][]Message{}
	}
	found := storage.Pool
	mutex.RUnlock()
	return found
}
