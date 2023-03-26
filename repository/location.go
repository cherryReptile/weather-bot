package repository

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"tgBotTask/domain"
	"time"
)

type LocationRepository interface {
	Create(loc *domain.Location) error
	Update(loc *domain.Location, username string, chatID uint) error
	Get(loc *domain.Location) error
	FindByChatID(loc *domain.Location, chatID uint)
}

type locationRepository struct {
	db *sqlx.DB
}

func NewLocationRepository(db *sqlx.DB) LocationRepository {
	return &locationRepository{
		db: db,
	}
}

func (c *locationRepository) Create(loc *domain.Location) error {
	loc.CreatedAt = time.Now()

	if _, err := c.db.NamedExec(`insert into chat_locations (username, lng, lat, weather_stat, country, city, chat_id, created_at, updated_at) values (:username, :lng, :lat, :weather_stat, :country, :city, :chat_id, :created_at, :updated_at)`, loc); err != nil {
		return err
	}

	if err := c.Get(loc); err != nil {
		return err
	}

	if loc.ID == 0 {
		return errors.New("record not found")
	}

	return nil
}

func (c *locationRepository) Update(loc *domain.Location, username string, chatID uint) error {
	now := time.Now()
	_, err := c.db.Exec("update chat_locations set username=$1, lng=$2, lat=$3, weather_stat=$4, country=$5, city=$6, updated_at=$7 where chat_id=$8", username, loc.Lng, loc.Lat, loc.WeatherStat, loc.Country, loc.City, now, chatID)
	if err != nil {
		return err
	}

	return nil
}

func (c *locationRepository) Get(loc *domain.Location) error {
	return c.db.Get(loc, "select * from chat_locations order by id desc limit 1")
}

func (c *locationRepository) FindByChatID(loc *domain.Location, chatID uint) {
	c.db.Get(loc, "select * from chat_locations where chat_id=$1 limit 1", chatID)
}
