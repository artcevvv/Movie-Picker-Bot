package main

import (
	"fmt"
	"strconv"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

// map for storing userState and userInputs

var userStates = make(map[int64]string)
var userInputs = make(map[int64]map[string]string)

// constants for states

const (
	stateWaitingForTitle = "waiting_for_title"
	stateWaitingForGenre = "waiting_for_genre"
	stateWaitingForBound = "waiting_for_bound"
)

func handleUserInput(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID
	username := update.Message.From.Username

	if state, exists := userStates[chatID]; exists {
		switch state {
		case stateWaitingForTitle:
			saveUserInput(chatID, "movieTitle", update.Message.Text)
			userStates[chatID] = stateWaitingForGenre
			_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "Enter the movie genre (or type 'skip' to leave it blank):"))

		case stateWaitingForGenre:
			genre := update.Message.Text
			if genre != "skip" {
				saveUserInput(chatID, "movieGenre", genre)
			}

			err := processMovieInput(username, chatID)

			if err != nil {
				_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), fmt.Sprintf("Failed to add movie: %v", err)))
			} else {
				_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "Movie added successfully!"))
			}

			delete(userStates, chatID)
			delete(userInputs, chatID)

		default:
			delete(userStates, chatID)
			delete(userInputs, chatID)
		}
	}
}

func saveUserInput(chatID int64, key, value string) {
	if userInputs[chatID] == nil {
		userInputs[chatID] = make(map[string]string)
	}
	userInputs[chatID][key] = value
}

func processMovieInput(username string, chatID int64) error {
	input := userInputs[chatID]

	telegramUsenameOwner := username
	telegramUserOwnerID := chatID

	movieTitle := input["movieTitle"]
	movieGenre := input["movieGenre"]

	var telegramUserBoundedID *int64

	if boundedIDText, exists := input["telegramUserBoundedID"]; exists && boundedIDText != "" {
		id, err := strconv.ParseInt(boundedIDText, 10, 64)

		if err != nil {
			return fmt.Errorf("invalid bounded ID: %v", err)
		}

		telegramUserBoundedID = &id
	}

	return addMovieHandler(telegramUsenameOwner, telegramUserOwnerID, movieTitle, movieGenre, telegramUserBoundedID)
}

// func processMovieRemoval(chatID int64, username, movieTitle string) error {
// 	input := userInputs[chatID]

// }
