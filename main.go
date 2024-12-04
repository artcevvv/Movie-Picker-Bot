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

	// commands
	botHandler.Handle(startCommand, th.CommandEqual("start"))
	botHandler.Handle(helpCommand, th.CommandEqual("help"))
	botHandler.Handle(stopCommand, th.CommandEqual("stop"))

	// movie commands
	botHandler.Handle(addMovie, th.CommandEqual("addmovie"))
	botHandler.Handle(getMovies, th.CommandEqual("getmovies"))
	botHandler.Handle(deleteMovieList, th.CommandEqual("deletemovie"))
	botHandler.Handle(randomCommand, th.CommandEqual("randmovie"))
	botHandler.Handle(getRandMovieByGenreHandler, th.CommandEqual("randbygenre"))

	// series commands
	botHandler.Handle(addSeries, th.CommandEqual("addseries"))
	botHandler.Handle(getSeries, th.CommandEqual("getseries"))
	botHandler.Handle(deleteSeries, th.CommandEqual("deleteseries"))
	botHandler.Handle(randomSeries, th.CommandEqual("randseries"))

	// dev & admin commands
	botHandler.Handle(todoCommand, th.CommandEqual("devtodo"))
	botHandler.Handle(sendGlobalAnnouncement, th.CommandEqual("announcement"))

	// CQ
	botHandler.HandleCallbackQuery(handleDeleteCallback, th.CallbackDataPrefix("delete:"))
	botHandler.HandleCallbackQuery(handleGenreSelect, th.CallbackDataPrefix("genre:"))
	botHandler.HandleCallbackQuery(handleSeriesGenreSelect, th.CallbackDataPrefix("seriesGenre:"))
	botHandler.HandleCallbackQuery(handleSeriesPaginationCQ, th.CallbackDataPrefix("seriesPage:"))
	botHandler.HandleCallbackQuery(handlePaginationCQ, th.CallbackDataPrefix("page:"))
	botHandler.HandleCallbackQuery(cqRandByGenre, th.CallbackDataPrefix("randbygenre:"))
	botHandler.HandleCallbackQuery(handleSeriesDeleteCQ, th.CallbackDataPrefix("deleteseries:"))

	// etc
	botHandler.Handle(anyText, th.AnyCommand())
	botHandler.Handle(handleUserInput)

	botHandler.Start()
}
