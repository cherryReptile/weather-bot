package repository

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"tgBotTask/domain"
	"time"
)

type historyRepository struct {
	db *sqlx.DB
}

func NewHistoryRepository(db *sqlx.DB) *historyRepository {
	return &historyRepository{
		db: db,
	}
}

func (c *historyRepository) Create(history *domain.History) error {
	history.CreatedAt = time.Now()

	if _, err := c.db.NamedExec(`insert into history (request, chat_id, created_at, updated_at) values (:request, :chat_id, :created_at, :updated_at)`, history); err != nil {
		return err
	}

	c.Get(history)
	if history.ID == 0 {
		return errors.New("history record not found")
	}

	return nil
}

func (c *historyRepository) Get(history *domain.History) {
	c.db.Get(history, "select * from history order by id desc limit 1")
}

func (c *historyRepository) GetAllByChatID(chatID uint) []*domain.History {
	var res []*domain.History
	c.db.Select(&res, "select * from history where chat_id=$1", chatID)
	return res
}

func (c *historyRepository) GetFirstRequest(history *domain.History, chatID uint) {
	c.db.Get(history, "select * from history where chat_id=$1 order by id limit 1", chatID)
}

func (c *historyRepository) GetLastRequest(history *domain.History, chatID uint) {
	c.db.Get(history, "select * from history where chat_id=$1 order by id desc limit 1", chatID)
}
