package registry

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/Peltoche/avro-gateway/internal"
)

// Client handle all the interaction between the service and the Schema .
type Client struct {
	client  *http.Client
	baseURL *url.URL
}

// NewClient instantiate a new Client.
func NewClient(schemaRegistryURL *url.URL) *Client {
	return &Client{
		client:  http.DefaultClient,
		baseURL: schemaRegistryURL,
	}
}

// FetchSchema corresponding to the subject/version.
func (t *Client) FetchSchema(ctx context.Context, subject string, version string) (string, error) {
	fetchSchemaPath, err := url.Parse(fmt.Sprintf("/subjects/%s/versions/%s/schema", subject, version))
	if err != nil {
		return "", internal.Errorf(internal.InternalError, "failed to generate the path: %s", err)
	}

	//nolint
	// Error not possible
	req, _ := http.NewRequest("GET", t.baseURL.ResolveReference(fetchSchemaPath).String(), nil)

	res, err := t.client.Do(req.WithContext(ctx))
	if err != nil {
		return "", internal.NewError(internal.RemoteError, err.Error())
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case 200:
		break
	case 404:
		return "", internal.Errorf(internal.NotFound, `schema %s/%s not found`, subject, version)
	default:
		return "", internal.Errorf(internal.RemoteError, "unexpected response status: %s", res.Status)
	}

	rawSchema, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", internal.Errorf(internal.RemoteError, "failed to read the response body: %s", err)
	}

	return string(rawSchema), nil
}
