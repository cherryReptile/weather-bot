package bootstrap

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"os"
)

type Database struct {
	Conn *sqlx.DB
}

func NewDbConn() *Database {
	return &Database{
		Conn: connectToDb(),
	}
}

func connectToDb() *sqlx.DB {
	db, err := sqlx.Connect("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	))

	if err != nil {
		logrus.Fatalf("Unable to connect to database: %v", err)
	}

	return db
}

func (d *Database) Close() {
	if err := d.Conn.Close(); err != nil {
		logrus.Fatalf("Unable to close database: %v", err)
		return
	}
}
