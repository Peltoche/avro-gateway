package schema

import (
	"context"
	"strconv"

	"github.com/Peltoche/avro-gateway/internal"
	"github.com/Peltoche/avro-gateway/model"
	uuid "github.com/satori/go.uuid"
)

// Usecase handling all the logic about the schema resource.
type Usecase struct {
	registry Registry
	storage  Storage
	// Set the uuid generation function as an attribute in order to be able to
	// mock id.
	generateUUID func() string
}

// Registry is used to fetch schema from any Schema Registry.
type Registry interface {
	FetchSchema(ctx context.Context, subject string, version string) (string, error)
}

// Storage used to persiste the clients state.
type Storage interface {
	RegisterNewClient(ctx context.Context, client *model.Client) error
}

// NewUsecase instantiate a new Usecase.
func NewUsecase(registry Registry, storage Storage) *Usecase {
	return &Usecase{
		registry: registry,
		storage:  storage,
		generateUUID: func() string {
			return uuid.NewV4().String()
		},
	}
}

// GetSchemaCmd is the requests parameters for the GetSchema method.
type GetSchemaCmd struct {
	Topic       string
	Application string
	Action      string
	Subject     string
	Version     string
}

// GetSchema check if the client is authorized to use the schema and return it.
func (t *Usecase) GetSchema(ctx context.Context, cmd *GetSchemaCmd) (string, error) {
	err := t.validateGetSchemaCmd(cmd)
	if err != nil {
		return "", err
	}

	schema, err := t.registry.FetchSchema(ctx, cmd.Subject, cmd.Version)
	if err != nil {
		return "", internal.Wrap(err, "failed to fetch the schema")
	}

	client := model.Client{
		ID:          t.generateUUID(),
		Topic:       cmd.Topic,
		Application: cmd.Application,
		Action:      cmd.Action,
		Subject:     cmd.Subject,
		Version:     cmd.Version,
	}

	err = t.storage.RegisterNewClient(ctx, &client)
	if err != nil {
		return "", internal.Wrap(err, "failed to register the client")
	}

	return schema, nil
}

func (t *Usecase) validateGetSchemaCmd(cmd *GetSchemaCmd) error {
	// Parse the "Version" field.
	if cmd.Version == "" {
		return internal.NewError(internal.ValidationError, `missing field "version"`)
	}
	if cmd.Version != "latest" {
		val, err := strconv.Atoi(cmd.Version)
		if err != nil || val < 1 {
			return internal.NewError(internal.ValidationError, `invalid input for field "version"`)
		}
	}

	// Parse the "Subject" field.
	if cmd.Subject == "" {
		return internal.NewError(internal.ValidationError, `missing field "subject"`)
	}

	// Parse the "Application" field.
	if cmd.Application == "" {
		return internal.NewError(internal.ValidationError, `missing field "application"`)
	}

	// Parse the "Topic" field.
	if cmd.Topic == "" {
		return internal.NewError(internal.ValidationError, `missing field "topic"`)
	}

	// Parse the "Action" field.
	if cmd.Action == "" {
		return internal.NewError(internal.ValidationError, `missing field "action"`)
	}
	if cmd.Action != "read" && cmd.Action != "write" {
		return internal.NewError(internal.ValidationError, `invalid input for field "action"`)
	}

	return nil
}
