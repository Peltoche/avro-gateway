package schema

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Peltoche/avro-gateway/internal"
	"github.com/gorilla/mux"
)

// HTTPHandler handling all the http logic about the schema resource.
type HTTPHandler struct {
	usecase usecase
}

type usecase interface {
	GetSchema(ctx context.Context, cmd *GetSchemaCmd) (string, error)
}

// NewHTTPHandler instantiate a new HTTPHandler.
func NewHTTPHandler(usecase usecase) *HTTPHandler {

	handler := &HTTPHandler{
		usecase: usecase,
	}

	return handler
}

// RegisterRoutes into the givem mux.Router.
func (t *HTTPHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/schema", t.Post).Methods("POST")
}

// Post /schemas/{subject}
func (t *HTTPHandler) Post(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Topic       string `json:"topic"`
		Application string `json:"application"`
		Version     string `json:"version"`
		Subject     string `json:"subject"`
		Action      string `json:"action"`
	}

	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		internal.WriteErrorIntoResponse(w, internal.NewError(internal.InvalidJSONBody, err.Error()))
		return
	}

	schema, err := t.usecase.GetSchema(r.Context(), &GetSchemaCmd{
		Topic:       req.Topic,
		Application: req.Application,
		Action:      req.Action,
		Subject:     req.Subject,
		Version:     req.Version,
	})
	if err != nil {
		internal.WriteErrorIntoResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(schema))
	if err != nil {
		log.Print(err)
	}
}
