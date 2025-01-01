package memory

func Delete(groupId string) {
	mutex.Lock()
	delete(storage.Pool, groupId)
	mutex.Unlock()
}
