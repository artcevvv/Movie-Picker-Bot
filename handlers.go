package main

import (
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

func stopCommand(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID
	delete(userStates, chatID)
	delete(userInputs, chatID)
	_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "Movie addition process stopped."))
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

	_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), movies))
}

func deleteMovie(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID
	username := update.Message.From.Username
	if username == "" {
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "You need a Telegram username to use this feature. Please set your username in Telegram settings."))
		return
	}

	if _, exists := userStates[chatID]; !exists {
		userStates[chatID] = stateWaitingForTitle
		userInputs[chatID] = make(map[string]string)
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "Enter the movie title:"))
	}
}

// func handleDeleteState(bot *telego.Bot, update telego.Update) {
// 	chatID := update.Message.Chat.ID
// 	username := update.Message.From.Username

// 	if state, exists := userStates[chatID]; exists {
// 		switch state {
// 		case stateWaitingForTitle:
// 			saveUserInput(chatID, "movieTitle", update.Message.Text)
// 			err := processMovieRemoval()
// 		}
// 	}
// }

// func addBound(bot *telego.Bot, update telego.Update) {
// 	chatID := update.Message.Chat.ID

// 	if _, exists := userStates[chatID]; !exists {
// 		userStates[chatID] = stateWaitingForBound
// 	}
// }

// case stateWaitingForBound:
// 	boundedIDText := update.Message.Text
// 	if boundedIDText != "skip" {
// 		saveUserInput(chatID, "telegramUserBoundedID", boundedIDText)
// 	}

// 	err := processMovieInput(chatID)
// 	if err != nil {
// 		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), fmt.Sprintf("Failed to add movie: %v", err)))
// 	} else {
// 		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "Movie added successfully!"))
// 	}

// 	delete(userStates, chatID)
// 	delete(userInputs, chatID)
