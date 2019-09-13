package telegram

import (
	"log"
	"strconv"

	"github.com/GwangGwang/ganeungbot/pkg/util"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

const tokenDir string = "/secrets/telegram"
const consoleChatIDDir string = "/secrets/telegram-consoleChatId"

// InitBot starts a bot
func InitBot(console chan string) {
	// TODO: move to util
	token := util.FileRead(tokenDir)
	chatIDStr := util.FileRead(consoleChatIDDir)
	consoleChatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	util.Check(err)

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	//bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for {
		select {
		case consoleMsg := <-console:
			msg := tgbotapi.NewMessage(consoleChatID, consoleMsg)
			sendThenLog(bot, &msg)
		case update := <-updates:
			if update.Message == nil {
				continue
			}
		}
	}

	//	for update := range updates {
	//		if update.Message == nil {
	//			continue
	//		}
	//
	//		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	//
	//		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	//		msg.ReplyToMessageID = update.Message.MessageID
	//
	//		bot.Send(msg)
	//	}
}

func sendThenLog(bot *tgbotapi.BotAPI, msg *tgbotapi.MessageConfig) {
	sendResult, err := bot.Send(msg)

	if err == nil {
		util.PrintChatLog(sendResult.Chat.ID, sendResult.MessageID, sendResult.From.UserName, sendResult.Text)
	}
}
