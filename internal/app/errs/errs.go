package errs

import "errors"

var (
	InvalidToken         = errors.New("invalid token")
	ResourceNotFound     = errors.New("resource not found")
	IncorrectCredentials = errors.New("incorrect credentials")
	Unauthorized         = errors.New("unauthorized")
)
