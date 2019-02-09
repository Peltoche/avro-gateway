package storage

import (
	"context"
	"testing"

	"github.com/Peltoche/avro-gateway/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_InMemory_RegisterNewClient_GetClientByID_success(t *testing.T) {
	client := model.Client{
		ID:          "some-id",
		Topic:       "some-topic",
		Application: "some-app",
		Action:      "read",
		Subject:     "my-avro-subject",
		Version:     "2",
	}

	storage := NewInMemory()

	err := storage.RegisterNewClient(context.Background(), &client)
	require.NoError(t, err)

	res, err := storage.GetClientByID(context.Background(), "some-id")

	require.NoError(t, err)
	assert.EqualValues(t, &client, res)
}

func Test_InMemory_GetClientByID_with_client_not_found(t *testing.T) {
	storage := NewInMemory()

	res, err := storage.GetClientByID(context.Background(), "some-unknown-id")

	require.NoError(t, err)
	assert.Nil(t, res)
}

func Test_InMemory_RegisterNewClient_twice(t *testing.T) {
	client := model.Client{
		ID:          "some-id",
		Topic:       "some-topic",
		Application: "some-app",
		Action:      "read",
		Subject:     "my-avro-subject",
		Version:     "2",
	}

	storage := NewInMemory()

	err := storage.RegisterNewClient(context.Background(), &client)
	require.NoError(t, err)

	err = storage.RegisterNewClient(context.Background(), &client)

	assert.EqualError(t, err, `internal error: storage conflict: try to register client "some-id" twice`)
}

func Test_InMemory_GetAllClientOnTopic_success(t *testing.T) {
	client := model.Client{
		ID:          "some-id",
		Topic:       "some-topic",
		Application: "some-app",
		Action:      "read",
		Subject:     "my-avro-subject",
		Version:     "2",
	}

	storage := NewInMemory()

	err := storage.RegisterNewClient(context.Background(), &client)
	require.NoError(t, err)

	res, err := storage.GetAllClientsOnTopic(context.Background(), "some-topic")

	require.NoError(t, err)
	assert.EqualValues(t, []model.Client{client}, res)
}
