package telegram

import (
	"io/ioutil"
	"log"
	"strings"
	"time"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

const tokenDir string = "/secrets/telegram"

// InitBot starts a bot
func InitBot() {
	data, err := ioutil.ReadFile(tokenDir)
	check(err)
	token := strings.TrimSpace(string(data))

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	// Optional: wait for updates and clear them if you don't want to handle
	// a large backlog of old messages
	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}

}
