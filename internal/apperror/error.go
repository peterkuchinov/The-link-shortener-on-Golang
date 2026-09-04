package apperror

import "errors"

var (
	ErrNotFound          = errors.New("short link not found")
	ErrCodeAlreadyExists = errors.New("custom code already exists")
	ErrInvalidCustomCode = errors.New("invalid custom code format")
)

type AppError struct {
	Err     error
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}
