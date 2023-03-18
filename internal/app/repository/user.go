package repository

import (
	"context"
	"fmt"
	"strings"

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

func (u *userRepository) FindAll(ctx context.Context, searchQuery dto.UserSearchQuery) (int, []domain.User, error) {
	var (
		fields = []string{}
		args   = []any{}
	)

	if searchQuery.Email != nil {
		fields = append(fields, " email ILIKE $1")
		args = append(args, "%"+*searchQuery.Email+"%")
	}

	var where string
	if len(fields) > 0 {
		where = " WHERE " + strings.Join(fields, " AND ")
	}
	totalQuery := fmt.Sprintf(`SELECT COUNT(*) AS total FROM users %s`, where)

	var total int
	err := u.DB.QueryRow(ctx, totalQuery, args...).Scan(&total)
	if err != nil {
		return 0, nil, errs.ClarifyError(err)
	}

	if total == 0 {
		return 0, nil, nil
	}

	sql := fmt.Sprintf(`SELECT * FROM users %s ORDER BY created_at DESC`, where)

	args = append(args, searchQuery.Limit)
	sql += fmt.Sprintf(` LIMIT $%d`, len(args))

	args = append(args, (searchQuery.Page-1)*searchQuery.Limit)
	sql += fmt.Sprintf(` OFFSET $%d`, len(args))

	rows, err := u.DB.Query(ctx, sql, args...)
	if err != nil {
		return 0, nil, errs.ClarifyError(err)
	}
	defer rows.Close()

	var users []domain.User
	users, err = pgx.CollectRows(rows, pgx.RowToStructByPos[domain.User])
	if err != nil {
		return 0, nil, errs.ClarifyError(err)
	}

	return total, users, nil
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
	err := row.Scan(
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
