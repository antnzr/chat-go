package errs

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	InvalidToken          = errors.New("invalid token")
	IncorrectCredentials  = errors.New("incorrect credentials")
	Unauthorized          = errors.New("unauthorized")
	InternalServerError   = errors.New("internal server error")
	ResourceNotFound      = errors.New("resource not found")
	ResourceAlreadyExists = errors.New("resource already exists")
	BadRequest            = errors.New("bad request")
	Forbidden             = errors.New("forbidden")
	DbError               = errors.New("db error")
	LimitExceeded         = errors.New("limit exceeded")
)

func ClarifyError(err error) error {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ResourceNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return ResourceAlreadyExists
		default:
			return fmt.Errorf("%q: %w", err, DbError)
		}
	}
	return err
}
