package bootstrap

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
	"os"
)

func NewBot() *telebot.Bot {
	pref := telebot.Settings{
		Token: os.Getenv("TG_BOT_TOKEN"),
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		logrus.Fatal(err)
		return nil
	}

	return b
}
