package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/minacio00/easyCourt/internal/model"
	"github.com/minacio00/easyCourt/internal/repository"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(user *model.User) error
	GetUserByID(id uint) (*model.User, error)
	GetAllUsers() ([]model.User, error)
	UpdateUser(user *model.User) error
	DeleteUser(id uint) error
	Authenticate(phone, password string) (*model.User, string, error)
}

type userService struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewUserService(repo repository.UserRepository) UserService {
	secret := viper.GetString("JWT_SECRET")
	return &userService{repo, secret}
}

func (s *userService) CreateUser(user *model.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return s.repo.CreateUser(user)
}

func (s *userService) GetUserByID(id uint) (*model.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *userService) GetAllUsers() ([]model.User, error) {
	return s.repo.GetAllUsers()
}

func (s *userService) UpdateUser(user *model.User) error {
	return s.repo.UpdateUser(user)
}

func (s *userService) DeleteUser(id uint) error {
	return s.repo.DeleteUser(id)
}

// need to return jwt
func (s *userService) Authenticate(phone, password string) (*model.User, string, error) {
	user, err := s.repo.GetUserByPhone(phone)
	if err != nil {
		return nil, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", err
	}
	tokenExpiresIn := 24 * time.Hour
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"phone":   user.Phone,
		"exp":     time.Now().Add(tokenExpiresIn).Unix(),
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, "", err
	}

	return user, tokenString, nil
}
