package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/dto"
	"github.com/antnzr/chat-go/internal/app/errs"
	"github.com/antnzr/chat-go/internal/app/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type userRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) domain.UserRepository {
	return &userRepository{
		DB: db,
	}
}

func (u *userRepository) Save(dto *dto.SignupRequest) (*domain.User, error) {
	conn, err := u.DB.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	sqlQuery := `
		INSERT INTO "users" ("email", "first_name", "last_name", "password")
		VALUES ($1, $2, $3, $4)
		RETURNING "id", "email", "password", "first_name", "last_name", "created_at";
	`
	row := conn.QueryRow(
		context.Background(),
		sqlQuery,
		&dto.Email,
		&dto.FirstName,
		&dto.LastName,
		&dto.Password,
	)

	user, err := scanRowsIntoUser(row)
	if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique") {
		return nil, errs.ResourceAlreadyExists
	} else if err != nil {
		return nil, errs.BadRequest
	}

	return user, nil
}

func (u *userRepository) FindByEmail(email string) (*domain.User, error) {
	conn, err := u.DB.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	sqlQuery := `
		SELECT "id", "email", "password", "first_name", "last_name", "created_at"
		FROM "users"
		WHERE "email" = $1;
	`
	row := conn.QueryRow(context.Background(), sqlQuery, email)

	user, err := scanRowsIntoUser(row)
	if err != nil {
		return nil, errs.ResourceNotFound
	}

	return user, nil
}

func (u *userRepository) FindAll() ([]domain.User, error) {
	return []domain.User{}, nil
}

func (u *userRepository) Delete(id int) error {
	ctx := context.Background()

	switch _, err := u.DB.Exec(ctx, `DELETE FROM "users" WHERE "id" = $1`, id); {
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return err
	case errors.Is(err, pgx.ErrNoRows):
		return errs.ResourceNotFound
	case err != nil:
		logger.Error("failed delete failed: %v\n", zap.Error(err))
		return errors.New("cannot delete product from database")
	}

	return nil
}

func (u *userRepository) FindById(id int) (*domain.User, error) {
	conn, err := u.DB.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	sqlQuery := `
		SELECT "id", "email", "password", "first_name", "last_name", "created_at"
		FROM "users"
		WHERE "id" = $1;
	`
	row := conn.QueryRow(context.Background(), sqlQuery, id)

	user, err := scanRowsIntoUser(row)
	if err != nil {
		return nil, errs.ResourceNotFound
	}

	return user, nil
}

func (u *userRepository) Update(userId int, dto *dto.UserUpdateRequest) (*domain.User, error) {
	conn, err := u.DB.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	const sqlQuery = `
		UPDATE "users" SET
			"first_name" = COALESCE(NULLIF($1,''), "first_name"),
			"last_name" = COALESCE(NULLIF($2,''), "last_name")
		WHERE "id" = $3
		RETURNING "id", "email", "password", "first_name", "last_name", "created_at";
	`
	row := conn.QueryRow(context.Background(), sqlQuery, dto.FirstName, dto.LastName, userId)
	user, err := scanRowsIntoUser(row)
	if err != nil {
		return nil, errs.ResourceNotFound
	}

	return user, nil
}

func scanRowsIntoUser(row pgx.Row) (*domain.User, error) {
	var user domain.User
	var err error
	err = row.Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
