package bootstrap

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
	"weatherBot/handlers"
)

func Start() {
	db := NewDbConn()
	bot := NewBot()
	NewServer(bot)

	m := handlers.NewMiddlewareHandler(db.Conn)
	handler := handlers.NewHandler(bot, db.Conn)

	bot.Use(middleware.Logger())
	bot.Use(m.HistoryMiddleware())

	bot.Handle("/start", handler.Start)
	bot.Handle("/location", handler.Location)
	bot.Handle(telebot.OnLocation, handler.HandleLocation)
	bot.Handle(telebot.OnEdited, handler.UpdateLocation)
	bot.Handle("/getweather", handler.GetWeather)
	bot.Handle("/getstats", handler.GetStats)

	logrus.Info("All handlers are set, starting...")
	bot.Start()
}
