package domain

import (
	"database/sql"
	"encoding/json"
	"time"
)

type History struct {
	ID        uint            `db:"id"`
	Request   json.RawMessage `db:"request"`
	ChatID    uint            `db:"chat_id"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt sql.NullTime    `db:"updated_at"`
}
