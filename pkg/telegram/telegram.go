// Client for Telegram messenger
package telegram

import (
	"log"

	"github.com/GwangGwang/ganeungbot/pkg/mid"
	"github.com/GwangGwang/ganeungbot/pkg/util"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// New instantiates a Telegram client and returns two channels, for receiving/sending
func New(token string) (chan mid.Msg, chan mid.Msg, error) {
	api, err := tgbotapi.NewBotAPI(token)
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
			formattedMsg := mid.Msg{
				Timestamp : int64(raw.Date),
				Username: raw.From.UserName,
				ChatID: raw.Chat.ID,
				Content: raw.Text,
			}
			receiveChan <- formattedMsg
		}
	}

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
