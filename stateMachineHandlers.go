package main

import (
	"fmt"
	"strconv"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

// map for storing userState and userInputs

// ["qwer", "Qwer", "qwer"]
// {"Key": "Value"}

var userStates = make(map[int64]string)
var userInputs = make(map[int64]map[string]string)

// constants for states

const (
	stateWaitingForTitle            = "waiting_for_title"
	stateWaitingForGenre            = "waiting_for_genre"
	stateWaitingForBoundUsername    = "waiting_for_bound"
	stateWaitingForPartnerAgreement = "waiting_for_partner_agreement"
)

func handleUserInput(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID

	if state, exists := userStates[chatID]; exists {
		switch state {
		case stateWaitingForTitle:
			saveUserInput(chatID, "movieTitle", update.Message.Text)
			userStates[chatID] = stateWaitingForGenre
			sendInitialGenreSelection(bot, chatID)
		case stateWaitingForGenre:
			// momma raised no bitch
		default:
			delete(userStates, chatID)
			delete(userInputs, chatID)
		}
	}
}

// func handleUserInputForBound(bot *telego.Bot, update telego.Update) {
// 	chatID := update.Message.Chat.ID

// 	if state, exists := userStates[chatID]; exists && state == stateWaitingForBoundUsername {
// 		switch state {
// 		case stateWaitingForBoundUsername:
// 			saveUserInput(chatID, "boundUsername", update.Message.Text)

// 		}
// 	}
// }

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

// initial message, which will be edited further
func sendInitialGenreSelection(bot *telego.Bot, chatID int64) (int, error) {
	itemsPerPage := 6

	start, end := 0, itemsPerPage
	if end > len(genres) {
		end = len(genres)
	}
	currentPageGenres := genres[start:end]

	var buttons [][]telego.InlineKeyboardButton
	var row []telego.InlineKeyboardButton

	for i, genre := range currentPageGenres {
		row = append(row, tu.InlineKeyboardButton(genre).WithCallbackData("genre:"+genre))

		if len(row) == 2 || i == len(currentPageGenres)-1 {
			buttons = append(buttons, row)
			row = nil
		}
	}

	var navRow []telego.InlineKeyboardButton
	if end < len(genres) {
		navRow = append(navRow, tu.InlineKeyboardButton("➡️ Next").WithCallbackData(fmt.Sprintf("page:%d", 1)))
	}
	buttons = append(buttons, navRow)

	buttons = append(buttons, tu.InlineKeyboardRow(tu.InlineKeyboardButton("Skip").WithCallbackData("genre:skip")))

	msg, err := bot.SendMessage(tu.Message(tu.ID(chatID), "Page 1:\n\nSelect a genre:").
		WithReplyMarkup(&telego.InlineKeyboardMarkup{InlineKeyboard: buttons}),
	)
	if err != nil {
		return 0, err
	}

	return msg.MessageID, nil
}

func editGenreSelection(bot *telego.Bot, chatID int64, messageID, page int) {
	itemsPerPage := 6

	start := page * itemsPerPage
	end := start + itemsPerPage
	if end > len(genres) {
		end = len(genres)
	}
	currentPageGenres := genres[start:end]

	var buttons [][]telego.InlineKeyboardButton
	var row []telego.InlineKeyboardButton

	for i, genre := range currentPageGenres {
		row = append(row, tu.InlineKeyboardButton(genre).WithCallbackData("genre:"+genre))

		if len(row) == 2 || i == len(currentPageGenres)-1 {
			buttons = append(buttons, row)
			row = nil
		}
	}

	var navRow []telego.InlineKeyboardButton
	if page > 0 {
		navRow = append(navRow, tu.InlineKeyboardButton("⬅️ Previous").WithCallbackData(fmt.Sprintf("page:%d", page-1)))
	}
	if end < len(genres) {
		navRow = append(navRow, tu.InlineKeyboardButton("➡️ Next").WithCallbackData(fmt.Sprintf("page:%d", page+1)))
	}
	buttons = append(buttons, navRow)

	buttons = append(buttons, tu.InlineKeyboardRow(tu.InlineKeyboardButton("Skip").WithCallbackData("genre:skip")))

	_, _ = bot.EditMessageText(&telego.EditMessageTextParams{
		ChatID:    tu.ID(chatID),
		MessageID: messageID,
		Text:      fmt.Sprintf("Page %d:\n\nPlease select a genre:", page+1),
		ReplyMarkup: &telego.InlineKeyboardMarkup{
			InlineKeyboard: buttons,
		},
	})
}
