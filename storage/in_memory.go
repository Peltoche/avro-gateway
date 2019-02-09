package storage

import (
	"context"

	"github.com/Peltoche/avro-gateway/model"
)

// InMemory storage without any persistence.
//
// It's mainly used for tests. It's not safe to use in production as it doesn't
// have any persistence!
type InMemory struct {
	clients map[string]model.Client
}

// NewInMemory instanciate a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{
		clients: map[string]model.Client{},
	}
}

// RegisterNewClient register a new Client into the list of clients.
func (t *InMemory) RegisterNewClient(ctx context.Context, client *model.Client) error {
	return nil
}
