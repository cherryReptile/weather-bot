package main

import (
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
	"tgBotTask/bootstrap"
	"tgBotTask/handlers"
)

func main() {
	db := bootstrap.NewDbConn()
	bot := bootstrap.NewBot()
	bot.Use(middleware.Logger())
	handler := handlers.NewHandler(bot, db.Conn)
	bot.Handle("/start", handler.Start)
	bot.Handle("/location", handler.Location)
	bot.Handle(telebot.OnLocation, handler.HandleLocation)
	bot.Handle(telebot.OnEdited, handler.UpdateLocation)
	bot.Handle("/getweather", handler.GetWeather)

	bot.Start()
}
