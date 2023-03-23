package repository

import (
	"github.com/jmoiron/sqlx"
	"tgBotTask/domain"
)

type ChatRepository interface {
	Create(chat *domain.Chat) error
}

type chatRepository struct {
	db *sqlx.DB
}

func NewChatRepository(db *sqlx.DB) ChatRepository {
	return &chatRepository{
		db: db,
	}
}

func (c *chatRepository) Create(chat *domain.Chat) error {
	return nil
}
