package schema

import (
	"context"
	"strconv"

	"github.com/Peltoche/schema-gateway/internal"
)

// Usecase handling all the logic about the schema resource.
type Usecase struct {
	registry Registry
}

// Registry is used to fetch schema from any Schema Registry.
type Registry interface {
	FetchSchema(ctx context.Context, subject string, version string) (string, error)
}

// NewUsecase instantiate a new Usecase.
func NewUsecase(registry Registry) *Usecase {
	return &Usecase{
		registry: registry,
	}
}

// GetSchemaCmd is the requests parameters for the GetSchema method.
type GetSchemaCmd struct {
	Action  string
	Subject string
	Version string
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

	// Parse the "Action" field.
	if cmd.Action == "" {
		return internal.NewError(internal.ValidationError, `missing field "action"`)
	}
	if cmd.Action != "read" && cmd.Action != "write" {
		return internal.NewError(internal.ValidationError, `invalid input for field "action"`)
	}

	return nil
}
