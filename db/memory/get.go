package memory

func Get(groupId string) []Message {
	mutex.RLock()
	found := storage[groupId]
	mutex.RUnlock()
	return found
}

func GetAll() map[string][]Message {
	mutex.RLock()
	found := storage
	mutex.RUnlock()
	return found
}
