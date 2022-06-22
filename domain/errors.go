package domain

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

var (
	// ErrRateLimitExceeded denotes an error raised when rate limit is exceeded
	ErrRateLimitExceeded = echo.NewHTTPError(http.StatusTooManyRequests, "rate limit exceeded")
	// ErrExtractorError denotes an error raised when extractor function is unsuccessful
	ErrExtractorError = echo.NewHTTPError(http.StatusForbidden, "error while extracting identifier")
)

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = errors.New("internal server error")
	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = errors.New("your request item not found")
	// ErrConflict will throw if the current action already exists
	ErrConflict = errors.New("your item already exists")
	// ErrBadParamInput will throw if the given request-body or params is not valid
	ErrBadParamInput = errors.New("given param is not valid")

	// ErrEmailAlreadyExists will throw if email already exists in system
	ErrEmailAlreadyExists = errors.New("email already exists")
	// ErrEmailOrPasswordNotMatch will throw if email or password not match in system
	ErrEmailOrPasswordNotMatch = errors.New("email or password not match")

	// ErrEmailNotExists will throw if email not exists in system
	ErrEmailNotExists = errors.New("email not exists")
	// ErrOTPWrongOrExpire will throw reset password
	ErrOTPWrongOrExpire = errors.New("otp wrong or expired")
)

func GetStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	log.Error(err)
	switch err {
	case ErrInternalServerError:
		return http.StatusInternalServerError
	case ErrNotFound:
		return http.StatusNotFound
	case ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
