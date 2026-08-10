package language

import (
	"sync"

	lang "cloud.google.com/go/language/apiv1"
	"github.com/google/uuid"
)

var (
	clients = make(map[string]*lang.Client)
	cMutex  = sync.RWMutex{}
)

func getClient(id string) (client *lang.Client) {
	cMutex.RLock()
	defer cMutex.RUnlock()

	return clients[id]
}
func addClient(client *lang.Client) string {
	cMutex.Lock()
	defer cMutex.Unlock()

	clientID := uuid.New().String()
	clients[clientID] = client
	return clientID
}

func removeClient(id string) {
	cMutex.Lock()
	defer cMutex.Unlock()

	delete(clients, id)
}
