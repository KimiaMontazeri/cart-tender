package model

import (
	"gorm.io/gorm"
	"time"
)

const (
	PENDING   = "PENDING"
	COMPLETED = "COMPLETED"
)

type (
	Cart struct {
		ID        int64 `gorm:"primary_key"`
		Username  string
		CreatedAt time.Time
		UpdatedAt time.Time
		Data      string
		State     string
	}

	CartRepo interface {
		Find(id int64) (Cart, error)
		FindByUser(username string) ([]Cart, error)
		FindAll() ([]Cart, error)
		Create(cart *Cart) error
		Update(cart *Cart) error
		Delete(id int64) error
	}

	SQLCartRepo struct {
		DB *gorm.DB
	}
)

func NewCart(username string, data string) *Cart {
	return &Cart{
		Username: username,
		Data:     data,
	}
}

func NewSQLCartRepo(db *gorm.DB) SQLCartRepo {
	return SQLCartRepo{DB: db}
}

func (r SQLCartRepo) Find(id int64) (Cart, error) {
	var stored Cart
	err := r.DB.Where(&Cart{ID: id}).First(&stored).Error

	return stored, err
}

func (r SQLCartRepo) FindAll() ([]Cart, error) {
	var result []Cart
	if err := r.DB.Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (r SQLCartRepo) FindByUser(username string) ([]Cart, error) {
	var result []Cart
	if err := r.DB.Find(&result).Where(&Cart{Username: username}).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (r SQLCartRepo) Create(cart *Cart) error {
	return r.DB.Create(cart).Error
}

func (r SQLCartRepo) Update(cart *Cart) error {
	cart.UpdatedAt = time.Now()
	return r.DB.Save(cart).Error
}

func (r SQLCartRepo) Delete(id int64) error {
	return r.DB.Delete(&Cart{ID: id}).Error
}
