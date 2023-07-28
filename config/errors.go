package config

import "fmt"

const (
	DATABASE_ERROR = "Database Error."
	SERVER_ERROR   = "Internal Server Error."
)

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (err AppError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s, Error: %v", err.Code, err.Message, err.Err)
}
