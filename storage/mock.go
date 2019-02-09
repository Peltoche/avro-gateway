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
