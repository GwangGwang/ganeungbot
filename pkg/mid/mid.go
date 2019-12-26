package mid

import (
	"fmt"
	"github.com/GwangGwang/ganeungbot/pkg/typehelper"
	"github.com/GwangGwang/ganeungbot/pkg/util"
	"github.com/GwangGwang/ganeungbot/pkg/weather"
	"log"
	"strings"
)

// Middleware is the middleware object that is the core of the bot
type Middleware struct {
	BotStartTime  int64
	ConsoleChatID int64
	ReceiveChan   chan Msg
	SendChan      chan Msg
	Weather       weather.Weather
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
	if m.ConsoleChatID != 0 {
		// Console for sending msgs via terminal
		go startConsole(m.SendChan, m.ConsoleChatID)
	}

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

		// Skip non-text msgs; TODO: future somehow support non-text msgs (image, sound, video, etc.)
		if msg.Content == "" {
			continue
		}

		// Skip stale msgs
		if msg.Timestamp < m.BotStartTime {
			log.Printf("msg %d, botstart %d", msg.Timestamp, m.BotStartTime)
			log.Printf("Skipping stale msg: %s", msg.Content)
			continue
		}

		// Print received msg in console
		util.PrintChatLog(msg.ChatID, 0, msg.Username, msg.Content)

		// Generate response based on msg text content
		responses := m.prepareResponse(msg.Username, msg.Content)

		for _, response := range responses {
			if len(response) > 0 {

				// DEBUG LINE
				//if msg.ChatID == -170492567 {
				//	continue
				//}
				//
				m.SendChan <- Msg{
					ChatID:   msg.ChatID,
					Content:  response,
					Username: "Ganeungbot",
				}
			}
		}
	}
}

// TODO: come up with a better function name
func (m *Middleware) prepareResponse(username string, txt string) []string {
	parseResult := Parse(txt)
	//log.Printf("%+v\n", parseResult)

	responses := []string{}
	if len(parseResult.Actions) > 0 {
		responses = m.buildResponse(parseResult.Actions[0], username, txt)
	}

	return responses
}

func (m *Middleware) buildResponse(action Action, username string, txt string) []string {
	resps := []string{}

	if answerList, ok := Answers[action]; ok {
		resps = append(resps, util.GetRandomElement(answerList))
	} else if action == ACTION_VERSUS {
		resps = append(resps, util.GetRandomElement(strings.Split(txt, "vs")))
		log.Printf("%+v\n", util.GetRandomElement(strings.Split(txt, "vs")))
	} else {
		switch action {
		case ACTION_TYPEHELPER:
			typehelpedMsg := typehelper.GetResponse(strings.Split(txt, typehelper.Trigger)[1])
			resps = append(resps, fmt.Sprintf("%s: %s", username, typehelpedMsg))

			// use the translated message for response generation
			resps = append(resps, m.prepareResponse(username, typehelpedMsg)...)

		case ACTION_WEATHER:
			// TODO: send in user's location or user's info so that we can fetch default location per user?
			resp, err := m.Weather.GetResponse(username, txt)
			if err != nil {
				log.Printf(err.Error())
				resps = append(resps, err.Error())
			} else {
				resps = append(resps, resp)
			}
		//	case ACTION_GAMESTATS:
		//		resp = chatObj.lolStats.GetResponse(txt)

		default:
		}
	}

	//fmt.Printf("%+v\n", resp)
	return resps
}
