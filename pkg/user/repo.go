package user

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNoUser  = errors.New("no user found")
	ErrBadPass = errors.New("invalid password")
)

type UserMemoryRepository struct {
	db *gorm.DB
}

func NewMemoryRepo(db *gorm.DB) *UserMemoryRepository {
	return &UserMemoryRepository{
		db: db,
	}
}

func (repo *UserMemoryRepository) Authorize(login, pass string) (*user, error) {
	users := make([]user, 0)
	repo.db.Where("login = ?", login).Find(&users)
	if len(users) != 1 {
		return nil, ErrNoUser
	}

	if users[0].Password != pass {
		return nil, ErrBadPass
	}

	return &users[0], nil
}

func (repo *UserMemoryRepository) AddUser(login, password string) error {
	newUser := user{
		Login:    login,
		Password: password,
	}

	db := repo.db.Create(&newUser)
	if db.Error != nil {
		return db.Error
	}

	return nil
}
