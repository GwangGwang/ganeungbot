package mid

import (
	"log"
	"strings"

	"github.com/GwangGwang/ganeungbot/pkg/util"
)

// Middleware is the middleware object that is the core of the bot
type Middleware struct {
	BotStartTime  int64
	ConsoleChatID int64
	ReceiveChan   chan Msg
	SendChan      chan Msg
}

// Msg is the received/sending message
type Msg struct {
	Timestamp int64
	ChatID    int64
	Username  string
	Content   string
}

// Start initializes the middlelayer for processing msgs received and to send
func (m *Middleware) Start() {
	// Console for sending msgs via terminal
	go startConsole(m.SendChan, m.ConsoleChatID)

	//var chats Chats = make(map[int64]Chat)

	m.process()
}

func (m *Middleware) process() {
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
			log.Printf("Skipping stale msg: %s", msg.Content)
			continue
		}

		util.PrintChatLog(msg.ChatID, 0, msg.Username, msg.Content)

		responseTxt := generateResponse(msg.Content)

		if len(responseTxt) > 0 {
			m.SendChan <- Msg{
				ChatID:   msg.ChatID,
				Content:  responseTxt,
				Username: "Ganeungbot",
			}
		}
	}
	//	msg := tgbotapi.NewMessage(t.ConsoleChatID, update.Message.Text)
	//	sendThenLog(bot, &msg)
}

func generateResponse(txt string) string {
	parseResult := Parse(txt)
	log.Printf("%+v\n", parseResult)

	var response string
	if len(parseResult.Actions) > 0 {
		response = buildResponse(parseResult.Actions[0], txt)
	}

	return response
}

func buildResponse(action Action, txt string) string {
	var resp string = ""

	if answerList, ok := Answers[action]; ok {
		resp = util.GetRandomElement(answerList)
	} else if action == ACTION_VERSUS {
		resp = util.GetRandomElement(strings.Split(txt, "vs"))
		log.Printf("%+v\n", util.GetRandomElement(strings.Split(txt, "vs")))
	} else {

		switch action {
		//	case ACTION_WEATHER:
		//		resp = weather.Process(text)

		//	case ACTION_GAMESTATS:
		//		resp = chatObj.lolStats.GetResponse(text)

		default:
		}
	}

	//fmt.Printf("%+v\n", resp)
	return resp
}
