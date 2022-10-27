package user

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNoUser  = errors.New("No user found")
	ErrBadPass = errors.New("Invald password")
)

type UserMemoryRepository struct {
	db   *gorm.DB
	data map[string]*user
}

func NewMemoryRepo(db *gorm.DB) *UserMemoryRepository {
	return &UserMemoryRepository{
		db: db,
		data: map[string]*user{
			"rvasily": &user{
				ID:       1,
				Login:    "rvasily",
				Password: "love",
			},
		},
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
