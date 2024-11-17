package main

import (
	"fmt"
	"os"

	"github.com/lpernett/godotenv"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func main() {
	InitDb()
	err := godotenv.Load()

	if err != nil {
		fmt.Println("ERR LOADING DOTENV FILE")
	}

	TELEGRAM_KEY := os.Getenv("TGKEY")

	bot, err := telego.NewBot(TELEGRAM_KEY, telego.WithDefaultDebugLogger())

	if err != nil {
		fmt.Println(err)
	}

	botUser, err := bot.GetMe()
	if err != nil {
		fmt.Println("ERR: ", err)
	}

	fmt.Printf("Bot user: %+v\n", botUser)

	updates, _ := bot.UpdatesViaLongPolling(nil)

	botHandler, _ := th.NewBotHandler(bot, updates)

	defer botHandler.Stop()
	defer bot.StopLongPolling()

	botHandler.Handle(startCommand, th.CommandEqual("start"))
	botHandler.Handle(helpCommand, th.CommandEqual("help"))
	botHandler.Handle(stopCommand, th.CommandEqual("stop"))
	botHandler.Handle(addMovie, th.CommandEqual("addmovie"))
	botHandler.Handle(getMovies, th.CommandEqual("getmovies"))

	botHandler.Handle(handleUserInput)

	botHandler.Start()
}
