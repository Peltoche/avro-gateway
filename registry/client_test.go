package registry

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Client_FetchSchema_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
		"subject":"foobar",
		"version":1,
		"id":1,
		"schema":"{ \"type\": \"record\", \"name\": \"Person\", \"namespace\": \"com.ippontech.kafkatutorials\", \"fields\": [ { \"name\": \"firstName\", \"type\": \"string\" }, { \"name\": \"lastName\", \"type\": \"string\" }, { \"name\": \"birthDate\", \"type\": \"long\" } ]}" }"
		}`))
	}))
	defer ts.Close()

	registryURL, err := url.Parse(ts.URL)
	require.NoError(t, err)

	client := NewClient(registryURL)

	schema, err := client.FetchSchema(context.Background(), "foobar", "1")

	require.NoError(t, err)
	assert.JSONEq(t, `{
	"type":"record",
	"name":"Person",
	"namespace":"com.ippontech.kafkatutorials",
	"fields":[
		{"name":"firstName","type":"string"},
		{"name":"lastName","type":"string"},
		{"name":"birthDate","type":"long"}
	]}`, schema)
}

func Test_Client_FetchSchema_with_path_error(t *testing.T) {
	registryURL, err := url.Parse("http://some-path")
	require.NoError(t, err)

	client := NewClient(registryURL)

	// Subject invalid in path
	schema, err := client.FetchSchema(context.Background(), "%gh&%ij", "1")

	assert.EqualError(t, err, "internal error: failed to generate the path: parse /subjects/%gh&%ij/versions/1: invalid URL escape \"%gh\"")
	assert.Empty(t, schema)
}

func Test_Client_FetchSchema_with_a_network_error(t *testing.T) {
	registryURL, err := url.Parse("invalid-url")
	require.NoError(t, err)

	client := NewClient(registryURL)

	// Subject invalid in path
	schema, err := client.FetchSchema(context.Background(), "foobar", "1")

	assert.EqualError(t, err, "remote error: Get /subjects/foobar/versions/1: unsupported protocol scheme \"\"")
	assert.Empty(t, schema)
}

func Test_Client_FetchSchema_whith_a_schema_not_found(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	registryURL, err := url.Parse(ts.URL)
	require.NoError(t, err)

	client := NewClient(registryURL)

	schema, err := client.FetchSchema(context.Background(), "foobar", "1")

	assert.Empty(t, schema)
	assert.EqualError(t, err, "not found: schema foobar/1 not found")
}

func Test_Client_FetchSchema_whith_an_unexpected_error(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))
	defer ts.Close()

	registryURL, err := url.Parse(ts.URL)
	require.NoError(t, err)

	client := NewClient(registryURL)

	schema, err := client.FetchSchema(context.Background(), "foobar", "1")

	assert.Empty(t, schema)
	assert.EqualError(t, err, "remote error: unexpected response status: 418 I'm a teapot")
}

func Test_Client_FetchSchema_with_an_invalid_response_body(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`not a json response`))
	}))
	defer ts.Close()

	registryURL, err := url.Parse(ts.URL)
	require.NoError(t, err)

	client := NewClient(registryURL)

	schema, err := client.FetchSchema(context.Background(), "foobar", "1")

	assert.Empty(t, schema)
	assert.EqualError(t, err, "remote error: invalid response body format: invalid character 'o' in literal null (expecting 'u')")
}
