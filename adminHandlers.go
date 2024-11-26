package main

import (
	"fmt"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func sendGlobalAnnouncement(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID

	isAdmin, err := checkIfAdmin(chatID)
	if err != nil {
		fmt.Printf("Something went wrong: %v", err)
	}

	if isAdmin {

		if err != nil {
			fmt.Printf("Failed to get unique chatid's: %v", err)
		}

		if _, exists := userStates[chatID]; !exists {
			userStates[chatID] = stateWaitingForAnnounceMsg
			userInputs[chatID] = make(map[string]string)
			_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "✒️ Write a message for global announcement:"))
		}
	} else {
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "You are not admin of the bot!"))
	}
}

func handleGlobalAnnouncementMsg(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID

	// Проверяем состояние пользователя
	if state, exists := userStates[chatID]; exists {
		switch state {
		case stateWaitingForAnnounceMsg:
			saveUserInput(chatID, "announceMessage", update.Message.Text)
			sendGlobalAnnouncementMsg(bot, chatID)
		default:
			// Если состояние не соответствует ожиданиям, очищаем данные
			delete(userStates, chatID)
			delete(userInputs, chatID)
		}
	}
}

func sendGlobalAnnouncementMsg(bot *telego.Bot, chatID int64) {
	globalAnnouncement = userInputs[chatID]["announceMessage"]
	_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), globalAnnouncement))
}
