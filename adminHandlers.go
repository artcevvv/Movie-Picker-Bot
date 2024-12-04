package main

import (
	"fmt"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func sendGlobalAnnouncement(bot *telego.Bot, update telego.Update) {
	liveChatID := update.Message.Chat.ID
	ChatIDs, err := getChatIDs()

	if err != nil {
		_, _ = bot.SendMessage(tu.Message(tu.ID(liveChatID), fmt.Sprintf("An error occured while sending global message: %v", err)))
		return
	}

	isAdmin, err := checkIfAdmin(liveChatID)
	if err != nil {
		fmt.Printf("Something went wrong: %v", err)
	}

	if isAdmin {

		if err != nil {
			fmt.Printf("Failed to get unique chatid's: %v", err)
		}

		for _, chatID := range ChatIDs {
			_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), globalAnnouncement).WithParseMode("HTML"))
		}
	} else {
		_, _ = bot.SendMessage(tu.Message(tu.ID(liveChatID), "You are not admin of the bot!"))
	}
}
