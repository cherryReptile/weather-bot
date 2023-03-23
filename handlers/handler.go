package handlers

import (
	"github.com/jmoiron/sqlx"
	"gopkg.in/telebot.v3"
	"tgBotTask/repository"
)

type Handler struct {
	Bot            *telebot.Bot
	userRepository repository.UserRepository
	chatRepository repository.ChatRepository
}

func NewHandler(bot *telebot.Bot, db *sqlx.DB) *Handler {
	return &Handler{
		Bot:            bot,
		userRepository: repository.NewUserRepository(db),
		chatRepository: repository.NewChatRepository(db),
	}
}

func (h *Handler) Start(c telebot.Context) error {
	return c.Send("Hello!")
}

func (h *Handler) GetWeather(c telebot.Context) error {
	return c.Send("Your weather: ")
}

func (h *Handler) GetStats(c telebot.Context) error {
	return c.Send("Your stats: ")
}
