package errs

import "errors"

var (
	InvalidToken          = errors.New("invalid token")
	IncorrectCredentials  = errors.New("incorrect credentials")
	Unauthorized          = errors.New("unauthorized")
	InternalServerError   = errors.New("internal server error")
	ResourceNotFound      = errors.New("resource not found")
	ResourceAlreadyExists = errors.New("resource already exists")
	BadRequest            = errors.New("bad request")
	Forbidden             = errors.New("forbidden")
)
