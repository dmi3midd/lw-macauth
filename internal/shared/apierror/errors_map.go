package apierror

import (
	"errors"

	"github.com/dmi3midd/lw-macauth/internal/service"
)

var ErrorMap = map[error]func(err error) error{
	service.ErrUnexpectedSigningMethod: func(err error) error {
		return NewInternalServerError(err)
	},
	service.ErrInvalidRefreshToken: func(err error) error {
		return NewBadRequestError(err, "Invalid refresh token")
	},
	service.ErrInvalidAccessToken: func(err error) error {
		return NewBadRequestError(err, "Invalid access token")
	},
	service.ErrSubjectAndIDNotFound: func(err error) error {
		return NewBadRequestError(err, "Invalid token")
	},
}

func MapError(err error) error {
	if err == nil {
		return nil
	}

	var apiErr APIError
	if errors.As(err, &apiErr) {
		return err
	}

	for serviceErr, mapFn := range ErrorMap {
		if errors.Is(err, serviceErr) {
			return mapFn(err)
		}
	}

	return NewInternalServerError(err)
}
