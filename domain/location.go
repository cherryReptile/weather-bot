package domain

import (
	"database/sql"
	"time"
)

type Location struct {
	ID          uint            `db:"id"`
	Username    string          `db:"username"`
	Lng         sql.NullFloat64 `db:"lng"`
	Lat         sql.NullFloat64 `db:"lat"`
	WeatherStat sql.NullString  `db:"weather_stat"`
	Country     sql.NullString  `db:"country"`
	City        sql.NullString  `db:"city"`
	ChatID      uint            `db:"chat_id"`
	CreatedAt   time.Time       `db:"created_at"`
	UpdatedAt   sql.NullTime    `db:"updated_at"`
}
