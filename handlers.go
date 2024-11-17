package main

import (
	"fmt"
	"math/rand"

	"github.com/mymmrac/telego"

	tu "github.com/mymmrac/telego/telegoutil"
)

func startCommand(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID

	_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), HelloWord))
}

func helpCommand(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID

	_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), helpMessage))
}

func todoCommand(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID

	_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), todoList))
}

func stopCommand(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID
	delete(userStates, chatID)
	delete(userInputs, chatID)
	_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "Movie addition process stopped."))
}

func getRandMovieByGenreHandler(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID
	username := update.Message.From.Username

	genres, err := getMoviesHandler(username)

	if err != nil {
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), fmt.Sprintf("Couldn't get genres! Error: %v", err)))
	}

	genreCount := make(map[string]int)

	for _, genre := range genres {
		movieGenre := genre["genre"]
		genreCount[movieGenre]++
	}

	var rows [][]telego.InlineKeyboardButton
	for genre, count := range genreCount {
		// Handle empty genre by setting a default label
		var genreLabel string
		if genre == "" {
			genreLabel = "No genre provided"
		} else {
			genreLabel = genre
		}

		// Create button with genre count
		button := tu.InlineKeyboardButton(fmt.Sprintf("%s (%d)", genreLabel, count)).WithCallbackData(fmt.Sprintf("randbygenre:%s", genre))

		rows = append(rows, []telego.InlineKeyboardButton{button})
	}

	_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "Select a genre:").WithReplyMarkup(&telego.InlineKeyboardMarkup{InlineKeyboard: rows}))
}

func randomCommand(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID
	username := update.Message.From.Username

	if update.Message.From.Username == "" {
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "You need a Telegram username to use this feature. Please set your username in Telegram settings."))
		return
	}

	movieList, err := getMoviesHandler(username)

	if err != nil {
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), fmt.Sprintf("Something went wrong! Error: %v", err)))
		return
	}

	if len(movieList) == 0 {
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "You have no movies in your list!"))
		return
	}

	randomIndex := rand.Intn(len(movieList))
	randomMovie := movieList[randomIndex]

	_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "Here is random movie from your list:\n\nTitle: <b>"+randomMovie["title"]+"</b>\nGenre: <b>"+randomMovie["genre"]+"</b>").WithParseMode("HTML"))
}

func addMovie(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID
	if update.Message.From.Username == "" {
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "You need a Telegram username to use this feature. Please set your username in Telegram settings."))
		return
	}

	if _, exists := userStates[chatID]; !exists {
		userStates[chatID] = stateWaitingForTitle
		userInputs[chatID] = make(map[string]string)
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "Enter the movie title:"))
	}
}

func getMovies(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID
	username := update.Message.From.Username
	if username == "" {
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "You need a Telegram username to use this feature. Please set your username in Telegram settings."))
		return
	}

	movies, err := getMoviesHandler(username)

	if err != nil {
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "No movies found for this user!\n\nTo add the movie, type /addmovie command"))
		return
	}

	var msg string

	for _, movie := range movies {
		if movie["genre"] == "" {
			msg += fmt.Sprintf("🎬 Title: %s\n", movie["title"])
		} else {
			msg += fmt.Sprintf("🎬 Title: %s | Genre: %s\n", movie["title"], movie["genre"])
		}
	}

	if msg == "" {
		msg = "No movies found for this user!"
	}

	_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "Here is movies you've added:\n\n"+msg))
}

func deleteMovieList(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID
	username := update.Message.From.Username
	if username == "" {
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "You need a Telegram username to use this feature. Please set your username in Telegram settings."))
		return
	}

	movies, err := getMoviesHandler(username)

	if err != nil {
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), fmt.Sprintf("Failed to fetch movies: %v", err)))
		return
	}

	var rows [][]telego.InlineKeyboardButton
	for _, movie := range movies {
		title := movie["title"]

		button := tu.InlineKeyboardButton(title).WithCallbackData(fmt.Sprintf("delete:%s", title))
		rows = append(rows, []telego.InlineKeyboardButton{button})
	}

	_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "Your movies:").WithReplyMarkup(&telego.InlineKeyboardMarkup{InlineKeyboard: rows}))
}
