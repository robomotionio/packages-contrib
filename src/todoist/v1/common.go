package todoist

import (
	"sync"

	"github.com/google/uuid"
)

var (
	clients = make(map[string]string)
	cMutex  = sync.RWMutex{}
)

func getToken(id string) (token string) {
	cMutex.RLock()
	defer cMutex.RUnlock()

	return clients[id]
}
func addToken(token string) string {
	cMutex.Lock()
	defer cMutex.Unlock()

	clientID := uuid.New().String()
	clients[clientID] = token
	return clientID
}

func removeToken(id string) {
	cMutex.Lock()
	defer cMutex.Unlock()

	delete(clients, id)
}
