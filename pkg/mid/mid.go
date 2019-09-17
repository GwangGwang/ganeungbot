// mid is the middlelayer for processing received msgs and preparing replies
package mid

import (
	"log"

	"github.com/GwangGwang/ganeungbot/pkg/util"


)

type Middleware struct {
	BotStartTime int64
	ConsoleChatID int64
	ReceiveChan chan Msg
	SendChan chan Msg
}

// Msg is the received/sending message
type Msg struct {
	Timestamp int64
	ChatID    int64
	Username  string
	Content string
}

// Start initializes the middlelayer for processing msgs received and to send
func (m *Middleware) Start() {
	// Console for sending msgs via terminal
	go startConsole(m.SendChan, m.ConsoleChatID)

	//var chats Chats = make(map[int64]Chat)

	for msg := range m.ReceiveChan {
		// Initialize new chat object if never seen before
	//	if chat, exists := chats[msg.ChatID]; !exists {
	//		chat = Chat{
	//			IsShutup: false,
	//		}
	//	}

		if msg.Content == "" {
			continue
		}

		if msg.Timestamp < m.BotStartTime {
			log.Printf("Not processing msg due to before bot start time")
			continue
		}

		util.PrintChatLog(msg.ChatID, 0, msg.Username, msg.Content)
	}
	//	msg := tgbotapi.NewMessage(t.ConsoleChatID, update.Message.Text)
	//	sendThenLog(bot, &msg)
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
