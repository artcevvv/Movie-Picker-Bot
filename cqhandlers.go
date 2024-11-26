package main

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func handleDeleteCallback(bot *telego.Bot, cq telego.CallbackQuery) {
	chatID := cq.Message.GetChat().ID
	username := cq.From.Username

	if username == "" {
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "You need a Telegram username to use this feature."))
		return
	}

	data := cq.Data
	if len(data) > 7 && data[:7] == "delete:" {
		movieTitle := data[7:]

		message, err := rmMovie(username, movieTitle)
		if err != nil {
			_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), fmt.Sprintf("Failed to remove movie: %v", err)))
			return
		}

		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), message))
	} else {
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "Invalid callback data."))
	}

	_ = bot.AnswerCallbackQuery(&telego.AnswerCallbackQueryParams{
		CallbackQueryID: cq.ID,
		Text:            "Processing your request...",
		ShowAlert:       false,
	})
}

func handleGenreSelect(bot *telego.Bot, cq telego.CallbackQuery) {
	chatID := cq.Message.GetChat().ID
	username := cq.Message.GetChat().Username

	if cq.Data == "" || len(cq.Data) < 7 || cq.Data[:6] != "genre:" {
		_ = bot.AnswerCallbackQuery(&telego.AnswerCallbackQueryParams{
			CallbackQueryID: cq.ID,
			Text:            "Invalid genre selection.",
		})
		return
	}

	selectedGenre := cq.Data[6:]

	_ = bot.AnswerCallbackQuery(&telego.AnswerCallbackQueryParams{
		CallbackQueryID: cq.ID,
		Text:            "Genre selected: " + selectedGenre,
	})

	saveUserInput(chatID, "movieGenre", selectedGenre)

	err := processMovieInput(username, chatID)

	if err != nil {
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), fmt.Sprintf("Failed to add movie: %v", err)))
	} else {
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "Movie added successfully!"))
	}

	delete(userStates, chatID)
	delete(userInputs, chatID)
}

func handleSeriesGenreSelect(bot *telego.Bot, cq telego.CallbackQuery) {
	chatID := cq.Message.GetChat().ID
	username := cq.Message.GetChat().Username

	// if cq.Data == "" || len(cq.Data) < 13 || cq.Data[:12] != "seriesGenre:" {
	// 	_ = bot.AnswerCallbackQuery(&telego.AnswerCallbackQueryParams{
	// 		CallbackQueryID: cq.ID,
	// 		Text:            "Invalid series genre selection.",
	// 	})
	// 	return
	// }

	selectedGenre := cq.Data[12:]

	_ = bot.AnswerCallbackQuery(&telego.AnswerCallbackQueryParams{
		CallbackQueryID: cq.ID,
		Text:            "Genre selected: " + selectedGenre,
	})

	saveUserInput(chatID, "seriesGenre", selectedGenre)

	err := processSeriesInput(username, chatID)

	if err != nil {
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), fmt.Sprintf("Failed to add series: %v", err)))
	} else {
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "Series added successfully!"))

	}

	delete(userStates, chatID)
	delete(userInputs, chatID)
}

func cqRandByGenre(bot *telego.Bot, cq telego.CallbackQuery) {
	chatID := cq.Message.GetChat().ID
	username := cq.From.Username

	if len(cq.Data) < 12 || cq.Data[:12] != "randbygenre:" || cq.Data == "" {
		_ = bot.AnswerCallbackQuery(&telego.AnswerCallbackQueryParams{
			CallbackQueryID: cq.ID,
			Text:            "Invalid genre selection.",
		})
		return
	}

	selectedGenre := cq.Data[12:]

	_ = bot.AnswerCallbackQuery(&telego.AnswerCallbackQueryParams{
		CallbackQueryID: cq.ID,
		Text:            "Genre selected: " + selectedGenre,
	})

	moviesByGenre, err := getMoviesByGenre(username, selectedGenre)

	if err != nil {
		_ = bot.AnswerCallbackQuery(&telego.AnswerCallbackQueryParams{
			CallbackQueryID: cq.ID,
			Text:            "You've selected incorrect genre",
		})
	} else {
		randomIndex := rand.Intn(len(moviesByGenre))
		randomVal := moviesByGenre[randomIndex]

		var genreLable string

		if selectedGenre == "" {
			genreLable = "No genre provided"
		} else {
			genreLable = selectedGenre
		}

		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "Here is random movie by genre <b>"+genreLable+"</b>: \n<b>"+randomVal["title"]+"</b>").WithParseMode("HTML"))
	}
}

// fucking pagination

func handlePaginationCQ(bot *telego.Bot, cq telego.CallbackQuery) {
	chatID := cq.Message.GetChat().ID
	messageID := cq.Message.GetMessageID()

	if len(cq.Data) > 5 && cq.Data[:5] == "page:" {
		page, err := strconv.Atoi(cq.Data[5:])
		if err != nil {
			_ = bot.AnswerCallbackQuery(&telego.AnswerCallbackQueryParams{
				CallbackQueryID: cq.ID,
				Text:            "Invalid page number.",
			})
			return
		}

		_ = bot.AnswerCallbackQuery(&telego.AnswerCallbackQueryParams{
			CallbackQueryID: cq.ID,
		})

		editGenreSelection(bot, chatID, messageID, page)
	}
}

func handleSeriesPaginationCQ(bot *telego.Bot, cq telego.CallbackQuery) {
	chatID := cq.Message.GetChat().ID
	messageID := cq.Message.GetMessageID()

	if len(cq.Data) > 11 && cq.Data[:11] == "seriesPage:" {
		page, err := strconv.Atoi(cq.Data[11:])
		if err != nil {
			_ = bot.AnswerCallbackQuery(&telego.AnswerCallbackQueryParams{
				CallbackQueryID: cq.ID,
				Text:            "Invalid page number",
			})
			return
		}

		_ = bot.AnswerCallbackQuery(&telego.AnswerCallbackQueryParams{
			CallbackQueryID: cq.ID,
		})

		editSeriesGenreSelection(bot, chatID, messageID, page)
	}
}
