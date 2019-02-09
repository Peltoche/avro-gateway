package storage

import (
	"context"
	"sync"

	"github.com/Peltoche/avro-gateway/internal"
	"github.com/Peltoche/avro-gateway/model"
)

// InMemory storage without any persistence.
//
// It's mainly used for tests. It's not safe to use in production as it doesn't
// have any persistence!
type InMemory struct {
	clients map[string]model.Client
	mutex   *sync.RWMutex
}

// NewInMemory instanciate a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{
		clients: map[string]model.Client{},
		mutex:   new(sync.RWMutex),
	}
}

// RegisterNewClient register a new Client into the list of clients.
func (t *InMemory) RegisterNewClient(ctx context.Context, client *model.Client) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	_, taken := t.clients[client.ID]
	if taken {
		return internal.Errorf(internal.InternalError, "storage conflict: try to register client %q twice", client.ID)
	}

	t.clients[client.ID] = *client

	return nil
}

// GetClientByID retrieve the client matching the id.
func (t *InMemory) GetClientByID(ctx context.Context, clientID string) (*model.Client, error) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	client, present := t.clients[clientID]
	if !present {
		return nil, nil
	}

	return &client, nil
}

// GetAllClientsOnTopic return all the client connected to a given topic.
func (t *InMemory) GetAllClientsOnTopic(ctx context.Context, topicName string) ([]model.Client, error) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	res := []model.Client{}
	for _, client := range t.clients {
		if client.Topic == topicName {
			res = append(res, client)
		}
	}

	return res, nil
}
