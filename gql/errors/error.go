package errors

import (
	"fmt"
	"net/http"
)

type ApiError struct {
	err        error
	ErrCode    string
	ErrMessage string
	HTTPCode   int
}

func (ae ApiError) Error() string {
	var err string
	if ae.err != nil {
		err = fmt.Sprintf("%s: %s", ae.ErrMessage, fmt.Errorf("%w", ae.err).Error())
	} else {
		err = fmt.Sprintf("%s", ae.ErrMessage)
	}

	return err
}

func (ae *ApiError) Unwrap() error {
	return ae.err
}

func NewApiError(err error, code string, msg string, httpCode int) ApiError {
	return ApiError{
		err:        err,
		ErrCode:    code,
		ErrMessage: msg,
		HTTPCode:   httpCode,
	}
}

func NewUnauthorizedError(err error) ApiError {
	return ApiError{
		err:        err,
		ErrCode:    "unauthorized",
		ErrMessage: "Unauthorized",
		HTTPCode:   http.StatusUnauthorized,
	}
}

func NewTOTPNotEnabledError(err error) ApiError {
	return ApiError{
		err:        err,
		ErrCode:    "totp-not-enabled",
		ErrMessage: "totp is not enabled for user",
		HTTPCode:   http.StatusBadRequest,
	}
}

func NewTOTPNotValidError(err error) ApiError {
	return ApiError{
		err:        err,
		ErrCode:    "totp-not-valid",
		ErrMessage: "TOTP not valid for user",
		HTTPCode:   http.StatusBadRequest,
	}
}
