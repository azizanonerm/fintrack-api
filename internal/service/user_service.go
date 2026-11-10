package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/azizanonerm/fintrack-api/internal/model"
	"github.com/azizanonerm/fintrack-api/internal/repository"
	"github.com/golang-jwt/jwt/v5"
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

type LoginInput struct {
	Email    string
	Password string
}

func (s *UserService) Login(ctx context.Context, input LoginInput, jwtSecret string) (string, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))

	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(input.Password),
	); err != nil {
		return "", errors.New("invalid email or password")
	}

	claims := jwt.MapClaims{
		"sub":   user.ID.String(),
		"email": user.Email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtSecret))
}
