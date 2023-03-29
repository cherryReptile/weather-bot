package storage

import "tgBotTask/domain"

type HistoryRepository interface {
	Create(history *domain.History) error
	GetAllByChatID(chatID uint) []*domain.History
	GetFirstRequest(history *domain.History, chatID uint)
	GetLastRequest(history *domain.History, chatID uint)
}
