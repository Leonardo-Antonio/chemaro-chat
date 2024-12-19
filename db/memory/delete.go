package memory

func Delete(groupId string) {
	mutex.Lock()
	delete(storage, groupId)
	mutex.Unlock()
}
