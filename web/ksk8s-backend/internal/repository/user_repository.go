package repository

import (
	"github.com/easzlab/ksk8s/internal/model"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Create(user *model.User) error {
	return DB.Create(user).Error
}

func (r *UserRepository) GetByID(id int64) (*model.User, error) {
	var user model.User
	err := DB.First(&user, id).Error
	return &user, err
}

func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *UserRepository) List() ([]model.User, error) {
	var users []model.User
	err := DB.Find(&users).Error
	return users, err
}

func (r *UserRepository) Update(user *model.User) error {
	return DB.Save(user).Error
}
