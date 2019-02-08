package schema

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// UsecaseMock is a mock implementation of schema.Usecase.
type UsecaseMock struct {
	mock.Mock
}

// GetSchema method mock.
func (t *UsecaseMock) GetSchema(ctx context.Context, cmd *GetSchemaCmd) (string, error) {
	args := t.Called(cmd)

	return args.String(0), args.Error(1)
}
