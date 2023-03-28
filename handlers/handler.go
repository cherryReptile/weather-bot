package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
	"tgBotTask/domain"
	"tgBotTask/pkg/weather"
	"tgBotTask/repository"
)

const (
	ErrMessage        = "Error, sorry. Try again..."
	AvailableCommands = "Available commands:\n\t\t/start: get commands list\n\t\t/getweather: your current weather if location is set\n\t\t/location: your current coordinates"
	NotFoundInDbErr   = "user not found in db"
	LocNotSet         = "Location didn't set. Please share it"
)

type Handler struct {
	Bot                *telebot.Bot
	locationRepository repository.LocationRepository
	historyRepository  repository.HistoryRepository
}

func NewHandler(bot *telebot.Bot, db *sqlx.DB) *Handler {
	return &Handler{
		Bot:                bot,
		locationRepository: repository.NewLocationRepository(db),
		historyRepository:  repository.NewHistoryRepository(db),
	}
}

func (h *Handler) Start(c telebot.Context) error {
	loc := new(domain.Location)
	h.locationRepository.FindByChatID(loc, uint(c.Chat().ID))
	if loc.ID == 0 {
		loc.ChatID = uint(c.Chat().ID)
		loc.Username = c.Chat().Username

		if err := h.locationRepository.Create(loc); err != nil {
			logrus.Error(err)
			return c.Send(ErrMessage)
		}
	}

	msg := fmt.Sprintf("Hello, %s! %s", c.Chat().Username, AvailableCommands)
	return c.Send(msg)
}

func (h *Handler) Location(c telebot.Context) error {
	loc := new(domain.Location)
	h.locationRepository.FindByChatID(loc, uint(c.Chat().ID))

	if loc.ID == 0 {
		logrus.Error(errors.New(NotFoundInDbErr))
		return c.Send(ErrMessage)
	}

	if !loc.Lng.Valid || !loc.Lat.Valid {
		return c.Send(LocNotSet)
	}

	msg := fmt.Sprintf("Your coordinates:\n\t\tLatitude: %f\n\t\tLongitude: %f", loc.Lat.Float64, loc.Lng.Float64)
	return c.Send(msg)
}

func (h *Handler) HandleLocation(c telebot.Context) error {
	loc := new(domain.Location)
	h.locationRepository.FindByChatID(loc, uint(c.Chat().ID))

	if loc.ID == 0 {
		logrus.Error(errors.New(NotFoundInDbErr))
		return c.Send(ErrMessage)
	}

	loc.Username = c.Chat().Username
	loc.Lng = sql.NullFloat64{Valid: true, Float64: float64(c.Message().Location.Lng)}
	loc.Lat = sql.NullFloat64{Valid: true, Float64: float64(c.Message().Location.Lat)}

	if err := h.locationRepository.Update(loc, c.Chat().Username, loc.ChatID); err != nil {
		logrus.Error(err)
		return c.Send(ErrMessage)
	}

	msg := fmt.Sprintf("Now you can use /getweather command")
	return c.Send(msg)
}

func (h *Handler) UpdateLocation(c telebot.Context) error {
	if c.Message().Location == nil {
		return nil
	}

	loc := new(domain.Location)
	h.locationRepository.FindByChatID(loc, uint(c.Chat().ID))

	if loc.ID == 0 {
		logrus.Error(errors.New(NotFoundInDbErr))
		return c.Send(ErrMessage)
	}

	if err := h.locationRepository.Update(loc, c.Chat().Username, loc.ChatID); err != nil {
		logrus.Error(err)
		return nil
	}

	return nil
}

func (h *Handler) GetWeather(c telebot.Context) error {
	loc := new(domain.Location)
	h.locationRepository.FindByChatID(loc, uint(c.Chat().ID))

	if loc.ID == 0 {
		logrus.Error(errors.New(NotFoundInDbErr))
		return c.Send(ErrMessage)
	}

	if !loc.Lng.Valid || !loc.Lat.Valid {
		return c.Send(LocNotSet)
	}

	w := new(weather.Info)
	err := w.GetWeather(loc.Lat.Float64, loc.Lng.Float64)

	if err != nil {
		logrus.Error(err)
		return c.Send("Failed to get weather stats, sorry 😔")
	}

	loc.WeatherStat = sql.NullString{Valid: true, String: w.Weather[0].Main}
	loc.Country = sql.NullString{Valid: true, String: w.Sys.Country}
	loc.City = sql.NullString{Valid: true, String: w.Name}

	if err = h.locationRepository.Update(loc, c.Chat().Username, loc.ChatID); err != nil {
		logrus.Error(err)
		return c.Send(ErrMessage)
	}

	return c.Send(w.GetInfo())
}

func (h *Handler) GetStats(c telebot.Context) error {
	var hh []*domain.History
	first := new(domain.History)
	last := new(domain.History)

	h.historyRepository.GetFirstRequest(first, uint(c.Chat().ID))
	if first.ID == 0 {
		logrus.Error(errors.New("first request of chat not found in history table"))
		return c.Send(ErrMessage)
	}

	h.historyRepository.GetLastRequest(last, uint(c.Chat().ID))
	if last.ID == 0 {
		logrus.Error(errors.New("last request of chat not found in history table"))
		return c.Send(ErrMessage)
	}

	fr := fmt.Sprintf("\n\t\tFirst request: %v", first.CreatedAt.Format("15:04:05 2006.01.02"))
	lr := fmt.Sprintf("\n\t\tLast request: %v", last.CreatedAt.Format("15:04:05 2006.01.02"))

	hh = h.historyRepository.GetAllByChatID(uint(c.Chat().ID))
	if len(hh) == 0 {
		return c.Send(ErrMessage)
	}

	msg := fmt.Sprintf("Your stats 🔢:%s%s\n\t\tTotal requests: %v", lr, fr, len(hh))
	return c.Send(msg)
}
