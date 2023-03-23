package domain

type Chat struct {
	ID     uint `json:"id" db:"id"`
	ChatID uint `json:"chat_id" db:"chat_id"`
}
