package apperror

import "errors"

var (
	ErrNotFound           = errors.New("short link not found")
	ErrCodeAlreadyExists  = errors.New("custom code already exists")
	ErrInvalidCustomCode  = errors.New("invalid custom code format")
	ErrInternal           = errors.New("internal server error")
)

type AppError struct {
	Err    error
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}
