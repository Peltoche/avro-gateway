package storage

import (
	"context"

	"github.com/Peltoche/avro-gateway/model"
	"github.com/stretchr/testify/mock"
)

// Mock implementation of a Storage.
type Mock struct {
	mock.Mock
}

// RegisterNewClient method mock.
func (t *Mock) RegisterNewClient(ctx context.Context, client *model.Client) error {
	return t.Called(client).Error(0)
}

// GetAllClientsOnTopic method mock.
func (t *Mock) GetAllClientsOnTopic(ctx context.Context, topicName string) ([]model.Client, error) {
	args := t.Called(topicName)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]model.Client), args.Error(1)
}
