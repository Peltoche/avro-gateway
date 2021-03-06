package schema

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Peltoche/avro-gateway/internal"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_HTTPHandler_Post_success(t *testing.T) {
	usecaseMock := new(UsecaseMock)

	handler := NewHTTPHandler(usecaseMock)

	usecaseMock.On("GetSchema", &GetSchemaCmd{
		Topic:       "my-topic",
		Application: "my-application",
		Action:      "read",
		Subject:     "my-avro-subject",
		Version:     "1",
	}).Return("some-schema", nil).Once()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "http://example.com/schema", strings.NewReader(`{
		"topic": "my-topic",
		"application": "my-application",
		"action": "read",
		"subject": "my-avro-subject",
		"version": "1"
	}`))

	router := mux.NewRouter()
	handler.RegisterRoutes(router)
	router.ServeHTTP(w, r)

	res := w.Result()
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "some-schema", string(body))

	usecaseMock.AssertExpectations(t)
}

func Test_HTTPHandler_Post_with_an_invalid_body_format(t *testing.T) {
	usecaseMock := new(UsecaseMock)

	handler := NewHTTPHandler(usecaseMock)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "http://example.com/schema", strings.NewReader("invalid json"))

	router := mux.NewRouter()
	handler.RegisterRoutes(router)
	router.ServeHTTP(w, r)

	res := w.Result()
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)
	assert.JSONEq(t, `{
		"kind": "invalid json body",
		"message": "invalid character 'i' looking for beginning of value"
	}`, string(body))

	usecaseMock.AssertExpectations(t)
}

func Test_HTTPHandler_Post_with_an_error_from_the_usecase(t *testing.T) {
	usecaseMock := new(UsecaseMock)

	handler := NewHTTPHandler(usecaseMock)

	usecaseMock.On("GetSchema", &GetSchemaCmd{
		Topic:       "my-topic",
		Application: "my-application",
		Action:      "read",
		Subject:     "my-avro-subject",
		Version:     "-1",
	}).Return("", internal.NewError(internal.ValidationError, "some-message")).Once()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "http://example.com/schema", strings.NewReader(`{
		"topic": "my-topic",
		"application": "my-application",
		"action": "read",
		"subject": "my-avro-subject",
		"version": "-1"
	}`))

	router := mux.NewRouter()
	handler.RegisterRoutes(router)
	router.ServeHTTP(w, r)

	res := w.Result()
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)
	assert.JSONEq(t, `{
		"kind": "validation error",
		"message": "some-message"
	}`, string(body))

	usecaseMock.AssertExpectations(t)
}

func Test_HTTPHandler_Post_with_an_unexpected_error_from_the_usecase(t *testing.T) {
	usecaseMock := new(UsecaseMock)

	handler := NewHTTPHandler(usecaseMock)

	usecaseMock.On("GetSchema", &GetSchemaCmd{
		Topic:       "my-topic",
		Application: "my-application",
		Action:      "read",
		Subject:     "my-avro-subject",
		Version:     "1",
	}).Return("", errors.New("some-unexpected-message")).Once()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "http://example.com/schema", strings.NewReader(`{
		"topic": "my-topic",
		"application": "my-application",
		"action": "read",
		"subject": "my-avro-subject",
		"version": "1"
	}`))

	router := mux.NewRouter()
	handler.RegisterRoutes(router)
	router.ServeHTTP(w, r)

	res := w.Result()
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.JSONEq(t, `{
		"kind": "internal error",
		"message": "unhandled error: some-unexpected-message"
	}`, string(body))

	usecaseMock.AssertExpectations(t)
}
