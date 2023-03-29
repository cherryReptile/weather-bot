package storage

import "tgBotTask/domain"

type LocationRepository interface {
	Create(loc *domain.Location) error
	Update(loc *domain.Location, username string, chatID uint) error
	FindByChatID(loc *domain.Location, chatID uint)
}
