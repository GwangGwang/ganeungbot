// mid is the middlelayer for processing received msgs and preparing replies
package mid

// Msg is the received/sending message
type Msg struct {
	Timestamp int64
	ChatID    string
	Username  string
	Content string
}

// Init initializes the middlelayer for processing msgs received and to send
func Init(receiveChan chan mid.Msg, sendChan chan mid.Msg, consoleChan chan mid.Msg) {
	go startConsole(sendChan)

}


	go func() {

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