package bootstrap

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
	"os"
	"time"
)

func NewBot() *telebot.Bot {
	pref := telebot.Settings{
		Token:  os.Getenv("TG_BOT_TOKEN"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		logrus.Fatal(err)
		return nil
	}

	return b
}
