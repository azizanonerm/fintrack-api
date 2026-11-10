package service

import (
	"context"
	"errors"
	"strings"

	"github.com/azizanonerm/fintrack-api/internal/model"
	"github.com/azizanonerm/fintrack-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository *repository.UserRepository
}

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) Register(ctx context.Context, input RegisterInput) (model.User, error) {
	name := strings.TrimSpace(input.Name)
	email := strings.ToLower(strings.TrimSpace(input.Email))

	if len(name) < 2 {
		return model.User{}, errors.New("name must be at least 2 characters")
	}

	if !strings.Contains(email, "@") {
		return model.User{}, errors.New("invalid email")
	}

	if len(input.Password) < 8 {
		return model.User{}, errors.New("password must be at least 8 characters")
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return model.User{}, err
	}

	return s.userRepository.Create(ctx, model.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(passwordHash),
	})
}
