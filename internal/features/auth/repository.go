package auth

import (
	"app/internal/database"

	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(user *database.User) error
	FindByEmail(email string) (*database.User, error)
	FindByID(id uint) (*database.User, error)
	FindByUsername(username string) (*database.User, error)
	UpdateUser(user *database.User) error
}

type authRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &authRepository{db: db}
}

func (r *authRepository) CreateUser(user *database.User) error {
	return r.db.Create(user).Error
}

func (r *authRepository) FindByEmail(email string) (*database.User, error) {
	var user database.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) FindByID(id uint) (*database.User, error) {
	var user database.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) FindByUsername(username string) (*database.User, error) {
	var user database.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) UpdateUser(user *database.User) error {
	return r.db.Save(user).Error
}
