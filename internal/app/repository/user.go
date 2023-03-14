package repository

import (
	"context"

	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/dto"
	"github.com/antnzr/chat-go/internal/app/errs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) domain.UserRepository {
	return &userRepository{
		DB: db,
	}
}

func (u *userRepository) Save(ctx context.Context, dto *dto.SignupRequest) (*domain.User, error) {
	sqlQuery := `
		INSERT INTO "users" ("email", "first_name", "last_name", "password")
		VALUES (@email, @firstName, @lastName, @password)
		RETURNING "id", "email", "password", "first_name", "last_name", "created_at";
	`
	args := pgx.NamedArgs{
		"email":     &dto.Email,
		"firstName": &dto.FirstName,
		"lastName":  &dto.LastName,
		"password":  &dto.Password,
	}
	row := u.DB.QueryRow(ctx, sqlQuery, args)

	user, err := scanRowsIntoUser(row)
	if err != nil {
		return nil, errs.ClarifyError(err)
	}

	return user, nil
}

func (u *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	sqlQuery := `
		SELECT "id", "email", "password", "first_name", "last_name", "created_at"
		FROM "users"
		WHERE "email" = $1;
	`
	row := u.DB.QueryRow(ctx, sqlQuery, email)

	user, err := scanRowsIntoUser(row)
	if err != nil {
		return nil, errs.ClarifyError(err)
	}

	return user, nil
}

func (u *userRepository) FindAll(ctx context.Context) ([]domain.User, error) {
	return []domain.User{}, nil
}

func (u *userRepository) Delete(ctx context.Context, id int) error {
	_, err := u.DB.Exec(ctx, `DELETE FROM "users" WHERE "id" = $1;`, id)
	if err != nil {
		return errs.ClarifyError(err)
	}
	return nil
}

func (u *userRepository) FindById(ctx context.Context, id int) (*domain.User, error) {
	sqlQuery := `
		SELECT "id", "email", "password", "first_name", "last_name", "created_at"
		FROM "users"
		WHERE "id" = $1;
	`
	row := u.DB.QueryRow(ctx, sqlQuery, id)

	user, err := scanRowsIntoUser(row)
	if err != nil {
		return nil, errs.ClarifyError(err)
	}

	return user, nil
}

func (u *userRepository) Update(ctx context.Context, userId int, dto *dto.UserUpdateRequest) (*domain.User, error) {
	args := pgx.NamedArgs{
		"firstName": &dto.FirstName,
		"lastName":  &dto.LastName,
		"id":        userId,
	}
	const sqlQuery = `
		UPDATE "users" SET
			"first_name" = COALESCE(NULLIF(@firstName, NULL), "first_name"),
			"last_name" = COALESCE(NULLIF(@lastName, NULL), "last_name")
		WHERE "id" = @id
		RETURNING "id", "email", "password", "first_name", "last_name", "created_at";
	`
	row := u.DB.QueryRow(ctx, sqlQuery, args)
	user, err := scanRowsIntoUser(row)
	if err != nil {
		return nil, errs.ClarifyError(err)
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
