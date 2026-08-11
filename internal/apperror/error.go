package apperror

import "errors"

var (
	ErrCodeAlreadyExists = errors.New("custom code is already taken")
	ErrInvalidCustomCode = errors.New("custom code contains invalid characters")
)
