package storage

import (
	"context"
	"testing"

	"github.com/Peltoche/avro-gateway/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_InMemory_RegisterNewClient_GetClientByID_success(t *testing.T) {
	storage := NewInMemory()

	client := model.Client{
		ID:          "some-id",
		Topic:       "some-topic",
		Application: "some-app",
		Action:      "read",
		Subject:     "my-avro-subject",
		Version:     "2",
	}

	err := storage.RegisterNewClient(context.Background(), &client)
	require.NoError(t, err)

	res, err := storage.GetClientByID(context.Background(), "some-id")

	require.NoError(t, err)
	assert.EqualValues(t, &client, res)
}

func Test_InMemory_RegisterNewClient_twice(t *testing.T) {
	storage := NewInMemory()

	client := model.Client{
		ID:          "some-id",
		Topic:       "some-topic",
		Application: "some-app",
		Action:      "read",
		Subject:     "my-avro-subject",
		Version:     "2",
	}

	err := storage.RegisterNewClient(context.Background(), &client)
	require.NoError(t, err)

	err = storage.RegisterNewClient(context.Background(), &client)

	assert.EqualError(t, err, `internal error: storage conflict: try to register client "some-id" twice`)
}
