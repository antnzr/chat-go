package repository

import (
	"context"
	"fmt"

	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/dto"
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

func (u *userRepository) Save(dto *dto.CreateUserRequest) (*domain.User, error) {
	conn, err := u.DB.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var user domain.User
	sqlQuery := `INSERT INTO users ("email", "first_name", "last_name", "password")
		VALUES ($1, $2, $3, $4)
		RETURNING "id", "email", "first_name", "last_name", "created_at";`
	row := conn.QueryRow(context.Background(), sqlQuery, &dto.Email, &dto.FirstName, &dto.LastName, &dto.Password)

	if err := row.Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName, &user.CreatedAt); err != nil {
		return &domain.User{}, err
	}
	fmt.Printf("user %v", &user)

	return &user, nil
}

func (u *userRepository) FindAll() ([]domain.User, error) {
	return []domain.User{}, nil
}

func (u *userRepository) Delete(id int) error {
	return nil
}

func (u *userRepository) FindById(id int) (*domain.User, error) {
	return &domain.User{}, nil
}
