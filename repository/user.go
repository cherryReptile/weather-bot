package repository

import (
	"github.com/jmoiron/sqlx"
	"tgBotTask/domain"
)

type UserRepository interface {
	Create(user *domain.User) error
	Get(user *domain.User)
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) Create(user *domain.User) error {
	return nil
}

func (u *userRepository) Get(user *domain.User) {

}
