package memory

func Set(groupId string, message Message) {
	mutex.Lock()
	storage[groupId] = append(storage[groupId], message)
	mutex.Unlock()
}
