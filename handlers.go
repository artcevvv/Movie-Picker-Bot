package main

import (
	"fmt"
	"math/rand"

	"github.com/mymmrac/telego"

	tu "github.com/mymmrac/telego/telegoutil"
)

// user -> addpartner -> username -> check if partner is user of bot ? bot sends message to partner "yesno" : err;

func anyText(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID

	_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), unknownCommand))
}

func startCommand(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID
	username := update.Message.From.Username

	kb := [][]telego.KeyboardButton{
		{
			tu.KeyboardButton("/addmovie"),
			tu.KeyboardButton("/getmovies"),
		},
		{
			tu.KeyboardButton("/randmovie"),
			tu.KeyboardButton("/randbygenre"),
		},
		// {
		// 	tu.KeyboardButton("/addpartner"),
		// 	tu.KeyboardButton("/partnerlist"),
		// },
		{
			tu.KeyboardButton("/suggest"),
		},
		{
			tu.KeyboardButton("/devtodo"),
		},
	}

	err := addUser(chatID, username)

	if err != nil {
		fmt.Printf("ERR ADDING USER: %v", err)
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "Something went wrong! Check your telegram username"))
	}

	_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), HelloWord).WithReplyMarkup(&telego.ReplyKeyboardMarkup{
		Keyboard:       kb,
		ResizeKeyboard: true,
		// OneTimeKeyboard: true,
	}))
}

func helpCommand(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID

	_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), helpMessage).WithParseMode("HTML"))
}

func todoCommand(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID

	_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), todoList).WithParseMode("HTML"))
}

func stopCommand(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID
	delete(userStates, chatID)
	delete(userInputs, chatID)
	_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), processStoppedMessage))
}

func getRandMovieByGenreHandler(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID
	username := update.Message.From.Username

	genres, err := getMoviesHandler(username)

	if err != nil {
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), fmt.Sprintf("⚠️ Couldn't get genres! Error: %v", err)))
	}

	genreCount := make(map[string]int)

	for _, genre := range genres {
		movieGenre := genre["genre"]
		genreCount[movieGenre]++
	}

	var rows [][]telego.InlineKeyboardButton
	for genre, count := range genreCount {
		var genreLabel string
		if genre == "" || genre == "skip" {
			genreLabel = "No genre provided"
		} else {
			genreLabel = genre
		}

		button := tu.InlineKeyboardButton(fmt.Sprintf("🎭 %s (%d)", genreLabel, count)).WithCallbackData(fmt.Sprintf("randbygenre:%s", genre))

		rows = append(rows, []telego.InlineKeyboardButton{button})
	}

	_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "🎥 Select a genre to proceed: 🍿").WithReplyMarkup(&telego.InlineKeyboardMarkup{InlineKeyboard: rows}))
}

func randomCommand(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID
	username := update.Message.From.Username

	if update.Message.From.Username == "" {
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), noUsername))
		return
	}

	movieList, err := getMoviesHandler(username)

	if err != nil {
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), fmt.Sprintf("⚠️ Something went wrong! Error: %v", err)))
		return
	}

	if len(movieList) == 0 {
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "You have no movies in your list! 😢 \n🎥 Add some with /addmovie to get started!"))
		return
	}

	randomIndex := rand.Intn(len(movieList))
	randomMovie := movieList[randomIndex]

	_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "🎰 Here is a random movie from your list:\n\n<b>🎬 Title:</b> "+randomMovie["title"]+"\n<b>🎭 Genre:</b> "+randomMovie["genre"]).WithParseMode("HTML"))
}

func addMovie(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID
	if update.Message.From.Username == "" {
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), noUsername))
		return
	}

	if _, exists := userStates[chatID]; !exists {
		userStates[chatID] = stateWaitingForTitle
		userInputs[chatID] = make(map[string]string)
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "✒️ Enter the movie title:"))
	}
}

func getMovies(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID
	username := update.Message.From.Username
	if username == "" {
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), noUsername))
		return
	}

	movies, err := getMoviesHandler(username)

	if err != nil {
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "You have no movies in your list! 😢 \n🎥 Add some with /addmovie to get started!"))
		return
	}

	var msg string

	for _, movie := range movies {
		if movie["genre"] == "" {
			msg += fmt.Sprintf("🎬 Title: %s\n", movie["title"])
		} else {
			msg += fmt.Sprintf("🎬 Title: %s | 🎭 Genre: %s\n", movie["title"], movie["genre"])
		}
	}

	if msg == "" {
		msg = "No movies found for this user!"
	}

	_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "📜 Here are the movies you've added:\n\n"+msg))
}

func deleteMovieList(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID
	username := update.Message.From.Username
	if username == "" {
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), noUsername))
		return
	}

	movies, err := getMoviesHandler(username)

	if err != nil {
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), fmt.Sprintf("❌ Failed to fetch movies: %v", err)))
		return
	}

	var rows [][]telego.InlineKeyboardButton
	for _, movie := range movies {
		title := movie["title"]

		button := tu.InlineKeyboardButton(title).WithCallbackData(fmt.Sprintf("delete:%s", title))
		rows = append(rows, []telego.InlineKeyboardButton{button})
	}

	_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "🗑️ Select which movie to delete:").WithReplyMarkup(&telego.InlineKeyboardMarkup{InlineKeyboard: rows}))
}

func addSeries(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID
	username := update.Message.From.Username
	if username == "" {
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), noUsername))
		return
	}

	if _, exists := userStates[chatID]; !exists {
		userStates[chatID] = stateWaitingForSeriesTitle
		userInputs[chatID] = make(map[string]string)
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "✒️ Enter the series title:"))
	}
}

func getSeries(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID
	listSeries, err := getSeriesHandler(chatID)

	if err != nil {
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), fmt.Sprintf("You have no series saved!\n\nTo add new series use /addseries command:\n\n %v", err)))
	}

	var msg string

	for _, series := range listSeries {
		if series["genre"] == "" {
			msg += fmt.Sprintf("🎬 Title: %s Number of episodes: %s\n", series["title"], series["episodes"])
		} else {
			msg += fmt.Sprintf("🎬 Title: %s Number of episodes: %s Genre: %s\n", series["title"], series["episodes"], series["genre"])
		}
	}

	if msg == "" {
		msg = "No series found for this user!"
	}

	_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "📜 Here are the series you've added:\n\n"+msg))
}

func deleteSeries(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID

	listSeries, err := getSeriesHandler(chatID)

	if err != nil {
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), fmt.Sprintf("❌ Failed to fetch series: %v", err)))
		return
	}

	var rows [][]telego.InlineKeyboardButton

	for _, series := range listSeries {
		title := series["title"]

		button := tu.InlineKeyboardButton(title).WithCallbackData(fmt.Sprintf("deleteseries:%s", title))
		rows = append(rows, []telego.InlineKeyboardButton{button})
	}
	_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "🗑️ Select which movie to delete:").WithReplyMarkup(&telego.InlineKeyboardMarkup{InlineKeyboard: rows}))
}
