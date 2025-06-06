package repository

import (
	"github.com/minacio00/easyCourt/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByID(id uint) (*model.User, error)
	GetUserByPhone(phone string) (*model.User, error)
	GetAllUsers(user *model.User) ([]model.User, error)
	UpdateUser(user *model.User) error
	DeleteUser(id uint) error
}
type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetUserByPhone(phone string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("phone = ?", phone).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetAllUsers(filter *model.User) ([]model.User, error) {
	var users []model.User
	query := r.db.Model(&model.User{})
	if filter != nil {
		conditions := make(map[string]interface{})

		if filter.Phone != "" {
			conditions["phone"] = filter.Phone
		}
		if filter.Email != "" {
			conditions["email"] = filter.Email
		}
		query.Where(conditions)
	}
	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) UpdateUser(user *model.User) error {
	existingUser := &model.User{}
	if err := r.db.First(existingUser, user.ID).Error; err != nil {
		return err // return error if user is not found
	}

	// Assign the found user's ID to the user object to ensure the correct record is updated
	user.ID = existingUser.ID

	// Update the user with the new data
	return r.db.Save(user).Error
}

func (r *userRepository) DeleteUser(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}
