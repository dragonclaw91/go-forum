package apperrs

import "errors"

// User-facing Errors
var (
	ErrUserNameTaken = errors.New("user name is not available")
	ErrInvalidPass   = errors.New("username or password does not match")
	ErrBlankFields   = errors.New("Fields can not be blank")
)

// generic catch all  error
var (
	Errgeneric = errors.New("The app is grumpy. Please try again later. ")
)

// Internal/System Errors
var (
	ErrDatabaseDown = errors.New("internal database connection error")
	ErrFailedToken  = errors.New("failed to generate secure token")
)
