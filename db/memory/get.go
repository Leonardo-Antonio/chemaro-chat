package memory

func Get(groupId string) []Message {
	mutex.RLock()
	found := storage[groupId]
	mutex.RUnlock()
	return found
}
