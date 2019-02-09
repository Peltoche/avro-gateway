package schema

import (
	"context"
	"errors"
	"testing"

	"github.com/Peltoche/avro-gateway/registry"
	"github.com/stretchr/testify/assert"
)

func Test_Usecase_GetSchema_success(t *testing.T) {
	registryMock := new(registry.Mock)

	usecase := NewUsecase(registryMock)

	registryMock.On("FetchSchema", "foobar", "1").Return("some-schema", nil).Once()

	schema, err := usecase.GetSchema(context.Background(), &GetSchemaCmd{
		Action:  "read",
		Subject: "foobar",
		Version: "1",
	})

	assert.NoError(t, err)
	assert.Equal(t, "some-schema", schema)

	registryMock.AssertExpectations(t)
}

func Test_Usecase_GetSchema_with_a_schema_validation_error(t *testing.T) {
	registryMock := new(registry.Mock)

	usecase := NewUsecase(registryMock)

	schema, err := usecase.GetSchema(context.Background(), &GetSchemaCmd{
		Action:  "read",
		Subject: "foobar",
		Version: "-1",
	})

	assert.EqualError(t, err, `validation error: invalid input for field "version"`)
	assert.Empty(t, schema)

	registryMock.AssertExpectations(t)
}

func Test_Usecase_GetSchema_with_a_fetch_schema_error(t *testing.T) {
	registryMock := new(registry.Mock)

	usecase := NewUsecase(registryMock)

	registryMock.On("FetchSchema", "foobar", "1").Return("", errors.New("some-error")).Once()

	schema, err := usecase.GetSchema(context.Background(), &GetSchemaCmd{
		Action:  "read",
		Subject: "foobar",
		Version: "1",
	})

	assert.EqualError(t, err, "internal error: failed to fetch the schema: some-error")
	assert.Empty(t, schema)

	registryMock.AssertExpectations(t)
}

func Test_Usecase_validateGetSchemaCmd(t *testing.T) {
	tests := []struct {
		Title string
		Cmd   GetSchemaCmd
		Err   string
	}{
		{
			Title: "valid",
			Cmd:   GetSchemaCmd{Action: "read", Subject: "bar", Version: "1"},
			Err:   "",
		},
		{
			Title: "missing_version",
			Cmd:   GetSchemaCmd{Action: "read", Subject: "bar", Version: ""},
			Err:   `validation error: missing field "version"`,
		},
		{
			Title: "version_is_latest",
			Cmd:   GetSchemaCmd{Action: "read", Subject: "bar", Version: "latest"},
			Err:   "",
		},
		{
			Title: "negative_version",
			Cmd:   GetSchemaCmd{Action: "read", Subject: "bar", Version: "-1"},
			Err:   `validation error: invalid input for field "version"`,
		},
		{
			Title: "invalid_version",
			Cmd:   GetSchemaCmd{Action: "read", Subject: "bar", Version: "foobar"},
			Err:   `validation error: invalid input for field "version"`,
		},
		{
			Title: "missing_subject",
			Cmd:   GetSchemaCmd{Action: "read", Subject: "", Version: "1"},
			Err:   `validation error: missing field "subject"`,
		},
		{
			Title: "missing_action",
			Cmd:   GetSchemaCmd{Action: "", Subject: "bar", Version: "1"},
			Err:   `validation error: missing field "action"`,
		},
		{
			Title: "action_read",
			Cmd:   GetSchemaCmd{Action: "read", Subject: "bar", Version: "1"},
			Err:   "",
		},
		{
			Title: "action_write",
			Cmd:   GetSchemaCmd{Action: "write", Subject: "bar", Version: "1"},
			Err:   "",
		},
		{
			Title: "invalid_action",
			Cmd:   GetSchemaCmd{Action: "invalid", Subject: "bar", Version: "1"},
			Err:   `validation error: invalid input for field "action"`,
		},
	}

	for _, test := range tests {
		t.Run(test.Title, func(tt *testing.T) {
			usecase := NewUsecase(nil)

			err := usecase.validateGetSchemaCmd(&test.Cmd)
			if test.Err == "" {
				assert.NoError(tt, err)
			} else {
				assert.EqualError(tt, err, test.Err)
			}
		})
	}
}
