package registry

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// Mock implementation of a Schema Registry.
type Mock struct {
	mock.Mock
}

// FetchSchema method mock.
func (t *Mock) FetchSchema(ctx context.Context, subject string, version string) (string, error) {
	args := t.Called(subject, version)

	return args.String(0), args.Error(1)
}
