package main

import (
	"fmt"
	"log"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func sendGlobalAnnouncement(bot *telego.Bot, update telego.Update) {
	liveChatID := update.Message.Chat.ID

	isAdmin, err := checkIfAdmin(liveChatID)
	if err != nil {
		fmt.Printf("Something went wrong: %v", err)
	}

	if isAdmin {
		chatIDs, err := getChatIDs()

		if err != nil {
			fmt.Printf("Failed to get unique chatid's: %v", err)
		}

		for _, chatID := range chatIDs {
			_, err := bot.SendMessage(tu.Message(tu.ID(chatID), globalAnnouncement))
			if err != nil {
				log.Printf("Failed to send message to chatID: %v", err)
			}
		}
	} else {
		_, _ = bot.SendMessage(tu.Message(tu.ID(liveChatID), "You are not admin of the bot!"))
	}
}

// func addAdmin(bot *telego.Bot, update telego.Update) {
// 	chatID := update.Message.Chat.ID

// 	isAdmin, err := checkIfAdmin(chatID)

// 	if err != nil {
// 		fmt.Printf("Something went wrong: %v", err)
// 	}

// 	if isAdmin {

// 	}

// }
