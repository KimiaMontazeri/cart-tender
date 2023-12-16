package model

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrUserExists = errors.New("username already exists")

type (
	User struct {
		Username string
		Password string
	}

	UserRepo interface {
		Find(username string) (User, error)
		FindAll() ([]User, error)
		Create(user *User) error
	}

	SQLUserRepo struct {
		DB *gorm.DB
	}
)

func NewUser(username string, password string) *User {
	return &User{
		Username: username,
		Password: password,
	}
}

func NewSQLUserRepo(db *gorm.DB) SQLUserRepo {
	return SQLUserRepo{DB: db}
}

func (r SQLUserRepo) Find(username string) (User, error) {
	var stored User
	err := r.DB.Where(&User{Username: username}).First(&stored).Error

	return stored, err
}

func (r SQLUserRepo) FindAll() ([]User, error) {
	var result []User
	if err := r.DB.Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (r SQLUserRepo) Create(user *User) error {
	_, err := r.Find(user.Username)
	if err == nil {
		return ErrUserExists
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashed)
	return r.DB.Create(user).Error
}

func (u *User) CheckPassword(p string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p))
	return err == nil
}
