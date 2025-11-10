package customerrors

import "fmt"

// AppError is a custom error type that includes an HTTP status code
type AppError struct {
	Code int
	Msg  string
}

// Error implements the built-in error interface
func (e *AppError) Error() string {
	return fmt.Sprintf("code %d: %s", e.Code, e.Msg)
}
