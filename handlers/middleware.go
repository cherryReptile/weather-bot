package handlers

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
	"sync"
	"weatherBot/domain"
	"weatherBot/storage/repository"
)

type MiddlewareHandler struct {
	historyCreator HistoryCreator
	mu             sync.Mutex
}

func NewMiddlewareHandler(db *sqlx.DB) *MiddlewareHandler {
	return &MiddlewareHandler{
		historyCreator: repository.NewHistoryRepository(db),
		mu:             sync.Mutex{},
	}
}

type HistoryCreator interface {
	Create(loc *domain.History) error
}

func (m *MiddlewareHandler) HistoryMiddleware() telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			go func() {
				history := new(domain.History)
				m.mu.Lock()
				body, err := json.Marshal(c.Update())
				m.mu.Unlock()
				if err != nil {
					logrus.Error(err)
					return
				}
				history.Request = body
				m.mu.Lock()
				history.ChatID = uint(c.Chat().ID)
				m.mu.Unlock()
				if err = m.historyCreator.Create(history); err != nil {
					logrus.Error(err)
					return
				}

				return
			}()
			return next(c)
		}
	}
}
