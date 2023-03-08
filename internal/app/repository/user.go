package repository

import (
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
	return &domain.User{
		Id:    1,
		Email: "ant@l.c",
	}, nil
}

func (u *userRepository) FindAll() ([]domain.User, error) {
	return []domain.User{}, nil
}

func (u *userRepository) Delete(id int) error {
	return nil
}

func (u *userRepository) FindById(id int) (*domain.User, error) {
	return &domain.User{
		Id:    1,
		Email: "ant@l.c",
	}, nil
}
