package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/dmi3midd/lw-macauth/internal/shared/apierror"
)

// Validator defines custom validation interface for request structures.
type Validator interface {
	Validate() error
}

// BindAndValidate decodes the request JSON body into the given type and validates it if it implements Validator.
func BindAndValidate[T any](r *http.Request) (T, error) {
	var body T
	if r.Body == nil {
		return body, apierror.NewBadRequestError(errors.New("request body is empty"), "Request body is empty")
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		if errors.Is(err, io.EOF) {
			return body, apierror.NewBadRequestError(err, "Request body is empty")
		}
		return body, apierror.NewBadRequestError(err, "Invalid request body")
	}
	defer r.Body.Close()

	if v, ok := any(body).(Validator); ok {
		if err := v.Validate(); err != nil {
			return body, apierror.NewBadRequestError(err, "Validation failed: "+err.Error())
		}
	} else if v, ok := any(&body).(Validator); ok {
		if err := v.Validate(); err != nil {
			return body, apierror.NewBadRequestError(err, "Validation failed: "+err.Error())
		}
	}

	return body, nil
}
