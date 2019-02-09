package schema

import (
	"context"
	"errors"
	"testing"

	"github.com/Peltoche/avro-gateway/internal"
	"github.com/Peltoche/avro-gateway/model"
	"github.com/Peltoche/avro-gateway/registry"
	"github.com/Peltoche/avro-gateway/storage"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Usecase_GetSchema_success(t *testing.T) {
	registryMock := new(registry.Mock)
	storageMock := new(storage.Mock)

	usecase := NewUsecase(registryMock, storageMock)
	usecase.generateUUID = func() string { return "some-id" }

	registryMock.On("FetchSchema", "foobar", "1").Return("some-schema", nil).Once()
	storageMock.On("GetAllClientsOnTopic", "some-topic").Return([]model.Client{}, nil).Once()
	storageMock.On("RegisterNewClient", &model.Client{
		ID:          "some-id",
		Topic:       "some-topic",
		Application: "my-application",
		Action:      "read",
		Subject:     "foobar",
		Version:     "1",
	}).Return(nil).Once()

	schema, err := usecase.GetSchema(context.Background(), &GetSchemaCmd{
		Topic:       "some-topic",
		Application: "my-application",
		Action:      "read",
		Subject:     "foobar",
		Version:     "1",
	})

	assert.NoError(t, err)
	assert.Equal(t, "some-schema", schema)

	registryMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Usecase_GetSchema_with_a_schema_validation_error(t *testing.T) {
	registryMock := new(registry.Mock)
	storageMock := new(storage.Mock)

	usecase := NewUsecase(registryMock, storageMock)

	schema, err := usecase.GetSchema(context.Background(), &GetSchemaCmd{
		Topic:       "some-topic",
		Application: "my-application",
		Action:      "read",
		Subject:     "foobar",
		Version:     "-1",
	})

	assert.EqualError(t, err, `validation error: invalid input for field "version"`)
	assert.Empty(t, schema)

	registryMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Usecase_GetSchema_with_a_fetch_schema_error(t *testing.T) {
	registryMock := new(registry.Mock)
	storageMock := new(storage.Mock)

	usecase := NewUsecase(registryMock, storageMock)

	registryMock.On("FetchSchema", "foobar", "1").Return("", errors.New("some-error")).Once()

	schema, err := usecase.GetSchema(context.Background(), &GetSchemaCmd{
		Topic:       "some-topic",
		Application: "my-application",
		Action:      "read",
		Subject:     "foobar",
		Version:     "1",
	})

	assert.EqualError(t, err, "internal error: failed to fetch the schema: some-error")
	assert.Empty(t, schema)

	registryMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Usecase_GetSchema_with_GetAllClientOnTopic_error(t *testing.T) {
	registryMock := new(registry.Mock)
	storageMock := new(storage.Mock)

	usecase := NewUsecase(registryMock, storageMock)
	usecase.generateUUID = func() string { return "some-id" }

	registryMock.On("FetchSchema", "foobar", "1").Return("some-schema", nil).Once()
	storageMock.On("GetAllClientsOnTopic", "some-topic").Return([]model.Client{}, errors.New("some-error")).Once()

	schema, err := usecase.GetSchema(context.Background(), &GetSchemaCmd{
		Topic:       "some-topic",
		Application: "my-application",
		Action:      "read",
		Subject:     "foobar",
		Version:     "1",
	})

	assert.Empty(t, schema)
	assert.EqualError(t, err, `internal error: failed to retrieve the list of clients connected to the topic "some-topic": some-error`)

	registryMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Usecase_GetSchema_with_in_incompatible_subject(t *testing.T) {
	registryMock := new(registry.Mock)
	storageMock := new(storage.Mock)

	usecase := NewUsecase(registryMock, storageMock)
	usecase.generateUUID = func() string { return "some-id" }

	registryMock.On("FetchSchema", "foobar", "1").Return("some-schema", nil).Once()
	storageMock.On("GetAllClientsOnTopic", "some-topic").Return([]model.Client{
		{
			ID: "some-other-id",
			// Same topic
			Topic:       "some-topic",
			Application: "an-other-application",
			Action:      "read",
			// Some other subject
			Subject: "an-other-subject",
			Version: "1",
		},
	}, nil).Once()

	schema, err := usecase.GetSchema(context.Background(), &GetSchemaCmd{
		Topic:       "some-topic",
		Application: "my-application",
		Action:      "read",
		Subject:     "foobar",
		Version:     "1",
	})

	assert.Empty(t, schema)
	assert.EqualError(t, err, `bad request: invalid subject: you can't use the subject "foobar" because the application "an-other-application" use the schema "an-other-subject/1"`)

	registryMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Usecase_GetSchema_with_a_register_client_error(t *testing.T) {
	registryMock := new(registry.Mock)
	storageMock := new(storage.Mock)

	usecase := NewUsecase(registryMock, storageMock)
	usecase.generateUUID = func() string { return "some-id" }

	registryMock.On("FetchSchema", "foobar", "1").Return("some-schema", nil).Once()
	storageMock.On("GetAllClientsOnTopic", "some-topic").Return([]model.Client{}, nil).Once()
	storageMock.On("RegisterNewClient", &model.Client{
		ID:          "some-id",
		Topic:       "some-topic",
		Application: "my-application",
		Action:      "read",
		Subject:     "foobar",
		Version:     "1",
	}).Return(internal.NewError(internal.InternalError, "some-error")).Once()

	schema, err := usecase.GetSchema(context.Background(), &GetSchemaCmd{
		Topic:       "some-topic",
		Application: "my-application",
		Action:      "read",
		Subject:     "foobar",
		Version:     "1",
	})

	assert.EqualError(t, err, "internal error: failed to register the client: some-error")
	assert.Empty(t, schema)

	registryMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Usecase_validateGetSchemaCmd(t *testing.T) {
	tests := []struct {
		Title string
		Cmd   GetSchemaCmd
		Err   string
	}{
		{
			Title: "valid",
			Cmd:   GetSchemaCmd{Topic: "some-topic", Application: "my-application", Action: "read", Subject: "bar", Version: "1"},
			Err:   "",
		},
		{
			Title: "missing_version",
			Cmd:   GetSchemaCmd{Topic: "some-topic", Application: "my-application", Action: "read", Subject: "bar", Version: ""},
			Err:   `validation error: missing field "version"`,
		},
		{
			Title: "version_is_latest",
			Cmd:   GetSchemaCmd{Topic: "some-topic", Application: "my-application", Action: "read", Subject: "bar", Version: "latest"},
			Err:   "",
		},
		{
			Title: "negative_version",
			Cmd:   GetSchemaCmd{Topic: "some-topic", Application: "my-application", Action: "read", Subject: "bar", Version: "-1"},
			Err:   `validation error: invalid input for field "version"`,
		},
		{
			Title: "invalid_version",
			Cmd:   GetSchemaCmd{Topic: "some-topic", Application: "my-application", Action: "read", Subject: "bar", Version: "foobar"},
			Err:   `validation error: invalid input for field "version"`,
		},
		{
			Title: "missing_subject",
			Cmd:   GetSchemaCmd{Topic: "some-topic", Application: "my-application", Action: "read", Subject: "", Version: "1"},
			Err:   `validation error: missing field "subject"`,
		},
		{
			Title: "missing_action",
			Cmd:   GetSchemaCmd{Topic: "some-topic", Application: "my-application", Action: "", Subject: "bar", Version: "1"},
			Err:   `validation error: missing field "action"`,
		},
		{
			Title: "action_read",
			Cmd:   GetSchemaCmd{Topic: "some-topic", Application: "my-application", Action: "read", Subject: "bar", Version: "1"},
			Err:   "",
		},
		{
			Title: "action_write",
			Cmd:   GetSchemaCmd{Topic: "some-topic", Application: "my-application", Action: "write", Subject: "bar", Version: "1"},
			Err:   "",
		},
		{
			Title: "invalid_action",
			Cmd:   GetSchemaCmd{Topic: "some-topic", Application: "my-application", Action: "invalid", Subject: "bar", Version: "1"},
			Err:   `validation error: invalid input for field "action"`,
		},
		{
			Title: "missing_topic",
			Cmd:   GetSchemaCmd{Topic: "", Application: "my-application", Action: "read", Subject: "bar", Version: "1"},
			Err:   `validation error: missing field "topic"`,
		},
		{
			Title: "missing_application",
			Cmd:   GetSchemaCmd{Topic: "some-topic", Application: "", Action: "read", Subject: "bar", Version: "1"},
			Err:   `validation error: missing field "application"`,
		},
	}

	for _, test := range tests {
		t.Run(test.Title, func(tt *testing.T) {
			usecase := NewUsecase(nil, nil)

			err := usecase.validateGetSchemaCmd(&test.Cmd)
			if test.Err == "" {
				assert.NoError(tt, err)
			} else {
				assert.EqualError(tt, err, test.Err)
			}
		})
	}
}

func Test_Usecase_generateUUID_is_a_valid_uuid(t *testing.T) {
	usecase := NewUsecase(nil, nil)

	res := usecase.generateUUID()

	assert.NotNil(t, uuid.FromStringOrNil(res))
}
