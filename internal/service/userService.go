package service

import (
	"errors"
	"fmt"
	"log"
	"strconv"
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
	Authenticate(identification, password string) (*model.User, string, string, error)
	RefreshToken(refreshToken string) (string, string, error)
	ValidateAccessToken(tokenstr string) (uint, error)
	IsUserAdmin(userID uint) (bool, error)
	ForgotPassword(user *model.User) error
}

type userService struct {
	repo      repository.UserRepository
	jwtSecret []byte
}

func NewUserService(repo repository.UserRepository) UserService {
	secret := viper.GetString("JWT_SECRET")
	return &userService{repo, []byte(secret)}
}

func (s *userService) ForgotPassword(user *model.User) error {
	if user.Email == "" {
		return errors.New("email não pode ser vazio")
	}
	return nil
	// check if the user exist, send token via email/sms
	// //maybe break this down into multiple methods: a service just for sending the token and other for checking and updatading the password
}

func (s *userService) ValidateAccessToken(tokenstr string) (uint, error) {
	token, err := jwt.Parse(tokenstr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signin method")
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return 0, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["type"] != "access" {
			return 0, errors.New("invalid token type")
		}
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return 0, errors.New("token expired")
			}
		} else {
			return 0, errors.New("invalid expiration claim")
		}
		if userID, ok := claims["user_id"].(float64); ok {
			return uint(userID), nil
		}
		return 0, errors.New("invalid user_id claim")
	}
	return 0, errors.New("invalid token")
}

func (s *userService) IsUserAdmin(userID uint) (bool, error) {
	user, err := s.GetUserByID(userID)
	if err != nil {
		return false, err
	}
	return user.IsAdmin, nil
}

func (s *userService) CreateUser(user *model.User) error {
	if err := user.Validate(); err != nil {
		return err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	usr, _ := s.repo.GetAllUsers(user)
	if len(usr) > 0 {
		log.Println(usr)
		return errors.New("essas credenciais já estão sendo utilizadas")
	}
	user.Password = string(hashedPassword)
	return s.repo.CreateUser(user)
}

func (s *userService) GetUserByID(id uint) (*model.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *userService) GetAllUsers() ([]model.User, error) {
	return s.repo.GetAllUsers(nil)
}

func (s *userService) UpdateUser(user *model.User) error {
	return s.repo.UpdateUser(user)
}

func (s *userService) DeleteUser(id uint) error {
	return s.repo.DeleteUser(id)
}

func (s *userService) Authenticate(identification, password string) (*model.User, string, string, error) {
	_, err := strconv.ParseInt(identification, 10, 0)
	filter := &model.User{}
	if err != nil {
		// if it's not possible to parser identifier as an int, assume the identitifier is an email
		filter.Email = identification
	} else {
		filter.Phone = identification
	}
	result, err := s.repo.GetAllUsers(filter)
	if err != nil {
		return nil, "", "", err
	}
	if len(result) > 1 {
		return nil, "", "", errors.New("erro ao buscar usuários")
	}
	if len(result) == 0 {
		return nil, "", "", errors.New("usuário não encontrado")
	}
	user := result[0]

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		fmt.Println(user.Password)
		fmt.Println("error during hashcompare")
		return nil, "", "", err
	}

	accessToken, err := s.generateToken(user.ID, "access", 15*time.Minute)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := s.generateToken(user.ID, "refresh", 7*24*time.Hour)
	if err != nil {
		return nil, "", "", err
	}

	return &user, accessToken, refreshToken, nil
}

func (s *userService) generateToken(userID uint, tokenType string, expiration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"type":    tokenType,
		"exp":     time.Now().Add(expiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *userService) RefreshToken(refreshToken string) (string, string, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return "", "", errors.New("invalid refresh token")
	}

	if claims["type"] != "refresh" {
		return "", "", errors.New("invalid token type")
	}

	userID := uint(claims["user_id"].(float64))

	newAccessToken, err := s.generateToken(userID, "access", 24*time.Hour)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := s.generateToken(userID, "refresh", 7*24*time.Hour)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}
