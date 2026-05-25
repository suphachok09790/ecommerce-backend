package service

import (
	"ecommerce/internal/model"
	"ecommerce/internal/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repository.UserRepository
}

func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(name, email, password string) error {
	// hash the password — never store plain text
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("could not hash password")
	}

	user := model.User{
		Name: name,
		Email: email,
		Password: string(hash),
	}
	return s.repo.Create(&user)
}

func (s *AuthService) Login(email, password string) (model.User, error) {
	// find user by email
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return model.User{} , errors.New("invalid credentials")
	}

	// compare password with stored hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return model.User{} , errors.New("invalid credentials")
	}

	return user, nil

}

