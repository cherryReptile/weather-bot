package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
	"sync"
	"tgBotTask/domain"
	"tgBotTask/pkg/weather"
	"tgBotTask/repository"
)

const (
	ErrMessage        = "Error, sorry. Try again..."
	AvailableCommands = "Available commands:\n\t\t/getweather: your current weather\n\t\t/location: your current coordinates"
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
		logrus.Error(errors.New("user not found in db"))
		return c.Send(ErrMessage)
	}

	if !loc.Lng.Valid || !loc.Lat.Valid {
		return c.Send(LocNotSet)
	}

	w := new(weather.WeatherResponse)
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
	return c.Send("Your stats: ")
}

func (h *Handler) HistoryMiddleware() telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			go func() {
				history := new(domain.History)
				mutex := sync.Mutex{}
				mutex.Lock()
				body, err := json.Marshal(c.Update())
				mutex.Unlock()
				if err != nil {
					logrus.Error(err)
					return
				}
				history.Request = body
				mutex.Lock()
				history.ChatID = uint(c.Chat().ID)
				mutex.Unlock()
				if err = h.historyRepository.Create(history); err != nil {
					logrus.Error(err)
					return
				}

				return
			}()
			return next(c)
		}
	}
}
