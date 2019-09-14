package telegram

import (
	"log"

	"github.com/GwangGwang/ganeungbot/pkg/util"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// Telegram is the service layer to Telegram chat
type Telegram struct {
	Token          string
	ConsoleChatID  int64
	UpdateChannel  tgbotapi.UpdatesChannel
	SendChannel    chan string
	ConsoleChannel chan string
}

// Start initiates chat
func (t *Telegram) Start() {

	bot, err := tgbotapi.NewBotAPI(t.Token)
	if err != nil {
		log.Panic(err)
	}

	//bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	t.UpdateChannel, err = bot.GetUpdatesChan(u)
	t.SendChannel = make(chan string)

	go func() {
		for consoleUpdate := range t.ConsoleChannel {
			msg := tgbotapi.NewMessage(t.ConsoleChatID, consoleUpdate)
			sendThenLog(bot, &msg)
		}
		// reaches here when channel closed
	}()

	for update := range t.UpdateChannel {
		if update.Message == nil {
			continue
		}
		msg := tgbotapi.NewMessage(t.ConsoleChatID, update.Message.Text)
		sendThenLog(bot, &msg)
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
