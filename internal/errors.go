package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// ErrorKind is the exostive list of all the possible type of error.
type ErrorKind string

var (

	// InternalError is returned in case of unexpected error.
	InternalError ErrorKind = "internal error"
	// RemoteError is returned in case of error from the Schema Registry.
	RemoteError ErrorKind = "remote error"
	// ValidationError returned in case of an invalid input.
	ValidationError ErrorKind = "validation error"
	// NotFound is returned if the schema is not found.
	NotFound ErrorKind = "not found"
	// InvalidJSONBody is returned when the server failed to decode the request body.
	InvalidJSONBody ErrorKind = "invalid json body"
	// BadRequest is returned when the client request is invalid.
	BadRequest ErrorKind = "bad request"
)

// Error returned a the differents components.
type Error struct {
	Kind    ErrorKind `json:"kind"`
	Message string    `json:"message"`
}

// Error is an implementation of error.
func (t *Error) Error() string {
	return fmt.Sprintf("%s: %s", t.Kind, t.Message)
}

// NewError instantiate a new error with the given type.
func NewError(errorKind ErrorKind, msg string) error {
	return &Error{
		Kind:    errorKind,
		Message: msg,
	}
}

// Errorf format the given message with the args then return an error with the given type.
func Errorf(errorKind ErrorKind, msg string, args ...interface{}) error {
	return &Error{
		Kind:    errorKind,
		Message: fmt.Sprintf(msg, args...),
	}
}

// Wrap the given message into a new error message
func Wrap(err error, msg string) error {
	innerError, ok := err.(*Error)
	if !ok {
		return &Error{
			Kind:    InternalError,
			Message: fmt.Sprintf("%s: %s", msg, err.Error()),
		}
	}

	return &Error{
		Kind:    innerError.Kind,
		Message: fmt.Sprintf("%s: %s", msg, innerError.Message),
	}
}

// Wrapf the given formated message into a new error message
func Wrapf(err error, msg string, args ...interface{}) error {
	formattedMsg := fmt.Sprintf(msg, args...)

	innerError, ok := err.(*Error)
	if !ok {
		return &Error{
			Kind:    InternalError,
			Message: fmt.Sprintf("%s: %s", formattedMsg, err.Error()),
		}
	}

	return &Error{
		Kind:    innerError.Kind,
		Message: fmt.Sprintf("%s: %s", formattedMsg, innerError.Message),
	}
}

// IsKind check if the given error is of the same kind of given errorKind.
func IsKind(errorKind ErrorKind, err error) bool {
	innerError, ok := err.(*Error)
	if !ok {
		return errorKind == InternalError
	}

	return innerError.Kind == errorKind

}

// WriteErrorIntoResponse will take the error and generate the correct response.
func WriteErrorIntoResponse(w http.ResponseWriter, err error) {
	innerError, ok := err.(*Error)
	if !ok {
		innerError = Errorf(InternalError, "unhandled error: %s", err).(*Error)
	}

	switch innerError.Kind {
	case RemoteError:
		w.WriteHeader(http.StatusBadGateway)
	case ValidationError:
		w.WriteHeader(http.StatusUnprocessableEntity)
	case InvalidJSONBody:
		w.WriteHeader(http.StatusUnprocessableEntity)
	case NotFound:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(innerError)
	if err != nil {
		log.Print(err)
	}
}
