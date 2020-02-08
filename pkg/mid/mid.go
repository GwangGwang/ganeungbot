package mid

import (
	"fmt"
	"github.com/GwangGwang/ganeungbot/pkg/translate"
	"github.com/GwangGwang/ganeungbot/pkg/typehelper"
	"github.com/GwangGwang/ganeungbot/pkg/util"
	"github.com/GwangGwang/ganeungbot/pkg/weather"
	"log"
	"strings"
)

// Middleware contains api destinations and information regarding the rooms the bot is invited to
type (
	Middleware struct {
		BotStartTime  int64
		ConsoleChatID int64
		ReceiveChan   chan Msg
		SendChan      chan Msg
		Weather       weather.Weather
		Translate     translate.Translate
		ChatGroups    map[int64]ChatGroup
	}

	ChatGroup struct {
		Users map[string]UserInfo
		IsShutup bool
	}

	// for now empty struct
	UserInfo struct {
	}

	// Msg is the received/sending message
	Msg struct {
		Timestamp int64
		ChatID    int64
		Username  string
		Content   string
	}
)

func New(startTime int64, receiveChan chan Msg, sendChan chan Msg, consoleChatId int64, w weather.Weather, t translate.Translate) *Middleware {
	return &Middleware{
		BotStartTime: startTime,
		ConsoleChatID: consoleChatId,
		ReceiveChan: receiveChan,
		SendChan: sendChan,
		Weather: w,
		Translate: t,
		ChatGroups: make(map[int64]ChatGroup),
	}
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
		if _, exists := m.ChatGroups[msg.ChatID]; !exists {
			m.ChatGroups[msg.ChatID] = ChatGroup{
				Users:                   make(map[string]UserInfo),
				IsShutup:                false,
			}
		}
		chatGroup := m.ChatGroups[msg.ChatID]

		// add user if not exists in chat group info yet
		// TODO: should probably skip over bots... how?
		if _, exists := chatGroup.Users[msg.Username]; !exists {
			chatGroup.Users[msg.Username] = UserInfo{}
		}

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
		responses := m.prepareResponse(msg)

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

		// update chatgroup
		m.ChatGroups[msg.ChatID] = chatGroup
	}
}

// TODO: come up with a better function name
func (m *Middleware) prepareResponse(msg Msg) []string {
	parseResult := Parse(msg.Content)
	//log.Printf("%+v\n", parseResult)

	responses := []string{}
	if len(parseResult.Actions) > 0 {
		responses = m.buildResponse(parseResult.Actions[0], msg)
	}

	return responses
}

func (m *Middleware) buildResponse(action Action, msg Msg) []string {
	username := msg.Username
	txt := msg.Content
	resps := []string{}

	if answerList, ok := Answers[action]; ok {
		resps = append(resps, util.GetRandomElement(answerList))
	} else if action == ACTION_VERSUS {
		resps = append(resps, util.GetRandomElement(strings.Split(txt, "vs")))
	} else {
		switch action {
		case ACTION_TRANSLATE:
			resp, err := m.Translate.GetResponse(msg.ChatID, txt)
			if err != nil {
				log.Printf(err.Error())
			}
			resps = append(resps, resp)
		case ACTION_TYPEHELPER:
			typehelpedMsg := typehelper.GetResponse(strings.Split(txt, typehelper.Trigger)[1])
			resps = append(resps, fmt.Sprintf("%s: %s", username, typehelpedMsg))

			// use the translated message for response generation
			resps = append(resps, m.prepareResponse(Msg{
				Timestamp: msg.Timestamp,
				ChatID:    msg.ChatID,
				Username:  username,
				Content:   typehelpedMsg,
			})...)
		case ACTION_WEATHER:
			// TODO: send in user's location or user's info so that we can fetch default location per user?
			resp, err := m.Weather.GetResponse(username, txt)

			// TODO: err should only be logged and the resp should contain err info as well (see translate)
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
