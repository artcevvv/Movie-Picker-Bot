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

// constants for states for movies

const (
	stateWaitingForTitle = "waiting_for_title"
	stateWaitingForGenre = "waiting_for_genre"
)

// constants for states for series

const (
	stateWaitingForSeriesTitle = "waiting_for_title"
	stateWaitingForSeriesGenre = "waiting_for_genre"
	stateWaitingForEpisodes    = "waiting_for_episodes"
)

// constants for states for user boundaries

// const (
// 	stateWaitingForBoundUsername    = "waiting_for_bound"
// 	stateWaitingForPartnerAgreement = "waiting_for_partner_agreement"
// )

// constants for states for admin actions

// const stateWaitingForAnnounceMsg = "waiting_for_announcement"

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

func handleUserSeriesAddition(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID

	if state, exists := userStates[chatID]; exists {
		switch state {
		case stateWaitingForSeriesTitle:
			saveUserInput(chatID, "seriesTitle", update.Message.Text)
			userStates[chatID] = stateWaitingForEpisodes
			_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "Enter number of episodes in series: "))
		case stateWaitingForEpisodes:
			saveUserInput(chatID, "seriesEpisodes", update.Message.Text)
			userStates[chatID] = stateWaitingForSeriesGenre
			sendInitialSeriesGenreSelection(bot, chatID)
		case stateWaitingForSeriesGenre:
			// momma raised no bitch
		default:
			delete(userStates, chatID)
			delete(userInputs, chatID)
		}
	}
}

func sendInitialSeriesGenreSelection(bot *telego.Bot, chatID int64) (int, error) {
	genresPerPage := 6

	start, end := 0, genresPerPage

	if end > len(seriesGenres) {
		end = len(genres)
	}

	currentPageGenres := seriesGenres[start:end]

	var buttons [][]telego.InlineKeyboardButton
	var row []telego.InlineKeyboardButton

	for i, genre := range currentPageGenres {
		row = append(row, tu.InlineKeyboardButton(genre).WithCallbackData("seriesGenre:"+genre))

		if len(row) == 2 || i == len(currentPageGenres)-1 {
			buttons = append(buttons, row)
			row = nil
		}
	}

	var navRow []telego.InlineKeyboardButton

	if end < len(seriesGenres) {
		navRow = append(navRow, tu.InlineKeyboardButton("➡️ Next").WithCallbackData(fmt.Sprintf("seriesPage:%d", 1)))
	}

	buttons = append(buttons, tu.InlineKeyboardRow(tu.InlineKeyboardButton("Skip").WithCallbackData("seriesGenre:skip")))

	buttons = append(buttons, navRow)

	msg, err := bot.SendMessage(tu.Message(tu.ID(chatID), "Page 1:\n\nSelect a genre:").
		WithReplyMarkup(&telego.InlineKeyboardMarkup{InlineKeyboard: buttons}),
	)

	if err != nil {
		return 0, err
	}

	return msg.MessageID, nil
}

func editSeriesGenreSelection(bot *telego.Bot, chatID int64, messageID, page int) {
	genresPerPage := 6

	start := page * genresPerPage
	end := start + genresPerPage

	if end > len(seriesGenres) {
		end = len(seriesGenres)
	}

	currentPageGenres := seriesGenres[start:end]

	var buttons [][]telego.InlineKeyboardButton
	var row []telego.InlineKeyboardButton

	for i, genre := range currentPageGenres {
		row = append(row, tu.InlineKeyboardButton(genre).WithCallbackData("seriesGenre:"+genre))

		if len(row) == 2 || i == len(currentPageGenres)-1 {
			buttons = append(buttons, row)
			row = nil
		}
	}

	var navRow []telego.InlineKeyboardButton

	if page > 0 {
		navRow = append(navRow, tu.InlineKeyboardButton("⬅️ Previous").WithCallbackData(fmt.Sprintf("seriesPage:%d", page-1)))
	}

	if end < len(seriesGenres) {
		navRow = append(navRow, tu.InlineKeyboardButton("➡️ Next").WithCallbackData(fmt.Sprintf("seriesPage:%d", 1)))
	}

	buttons = append(buttons, tu.InlineKeyboardRow(tu.InlineKeyboardButton("Skip").WithCallbackData("seriesGenre:skip")))

	buttons = append(buttons, navRow)

	_, _ = bot.EditMessageText(&telego.EditMessageTextParams{
		ChatID:    tu.ID(chatID),
		MessageID: messageID,
		Text:      fmt.Sprintf("Page %d:\n\nSelect a genre:", page+1),
		ReplyMarkup: &telego.InlineKeyboardMarkup{
			InlineKeyboard: buttons,
		},
	})
}

func processSeriesInput(username string, chatID int64) error {
	input := userInputs[chatID]

	telegramUsername := username
	telegramUserID := chatID

	seriesTitle := input["seriesTitle"]
	seriesEpisodes := input["seriesEpisodes"]
	seriesGenre := input["seriesGenre"]

	return addSeriesHandler(telegramUsername, telegramUserID, seriesTitle, seriesEpisodes, seriesGenre)
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
