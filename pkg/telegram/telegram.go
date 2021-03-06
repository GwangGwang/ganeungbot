package telegram

import (
	"fmt"
	"github.com/GwangGwang/ganeungbot/pkg/mid"
	"log"

	"github.com/GwangGwang/ganeungbot/pkg/util"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// New instantiates a Telegram client and returns two channels, for receiving/sending
func New(apiKey string) (chan mid.Msg, chan mid.Msg, error) {
	log.Println("Initializing telegram pkg")

	if len(apiKey) == 0 {
		return nil, nil, fmt.Errorf("no telegram api key supplied")
	}

	// Start bot api
	api, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		return nil, nil, err
	}

	//bot.Debug = true
	log.Printf("Authorized on account %s", api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Main channel for receiving msg from Telegram
	rawReceiveChan, err := api.GetUpdatesChan(u)
	receiveChan := make(chan mid.Msg)

	// Need to reformat Telegram-specific msg structure to our own
	go func() {
		for raw := range rawReceiveChan {

			// TODO: why NPE?
			if raw.Message == nil {
				continue
			}
			log.Printf("%+v\n", raw)
			formattedMsg := mid.Msg{
				Timestamp: int64(raw.Message.Date),
				Username:  raw.Message.From.UserName,
				ChatID:    raw.Message.Chat.ID,
				Content:   raw.Message.Text,
			}
			receiveChan <- formattedMsg
		}
	}()

	sendChan := make(chan mid.Msg)
	go listenOutgoing(api, sendChan)

	return receiveChan, sendChan, nil
}

// listenOutgoing enables Telegram API to process msgs incoming to channel and send to API
func listenOutgoing(api *tgbotapi.BotAPI, sendChan chan mid.Msg) {
	for msgToSend := range sendChan {
		// TODO: error checking?
		outgoingMsg := tgbotapi.NewMessage(msgToSend.ChatID, msgToSend.Content)
		sendResult, err := api.Send(outgoingMsg)
		if err == nil {
			util.PrintChatLog(sendResult.Chat.ID, sendResult.MessageID, sendResult.From.UserName, sendResult.Text)
		}
	}
}
