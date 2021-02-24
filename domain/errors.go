package domain

import "errors"

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = errors.New("Internal Server Error")
	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = errors.New("Your requested Item is not found")
	// ErrConflict will throw if the current action already exists
	ErrConflict = errors.New("Your Item already exist")
	// ErrBadParamInput will throw if the given request-body or params is not valid
	ErrBadParamInput = errors.New("Given Param is not valid")

	// ErrEmailAlreadyExists will throw if email already exists in system
	ErrEmailAlreadyExists = errors.New("Email already exists")
	// ErrEmailOrPasswordNotMatch will throw if email or password not match in system
	ErrEmailOrPasswordNotMatch = errors.New("Email or password not match")

	// ErrEmailNotExists will throw if email not exists in system
	ErrEmailNotExists = errors.New("Email not exists")
	// ErrOTPWrongOrExpire will throw reset password
	ErrOTPWrongOrExpire = errors.New("OTP wrong or expired")
)
