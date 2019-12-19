package mid

import (
	"github.com/GwangGwang/ganeungbot/pkg/typehelper"
	"log"
	"strconv"
	"strings"

	"github.com/GwangGwang/ganeungbot/internal/pkg/config"
	"github.com/GwangGwang/ganeungbot/pkg/util"
	"github.com/GwangGwang/ganeungbot/pkg/weather"
)

// Middleware is the middleware object that is the core of the bot
type Middleware struct {
	BotStartTime  int64
	ConsoleChatID int64
	ReceiveChan   chan Msg
	SendChan      chan Msg
	Weather       weather.Instance
}

// Msg is the received/sending message
type Msg struct {
	Timestamp int64
	ChatID    int64
	Username  string
	Content   string
}

const consoleChatIDdir = "consoleChatID"

// New initializes and returns a new middleware instance
func New(startTime int64, receiveChan chan Msg, sendChan chan Msg, w weather.Instance) Middleware {
	log.Println("Initializing middleware")

	consoleChatID, err := readConfig()
	if err != nil || consoleChatID == 0 {
		log.Printf("Problem when fetching console chat id: %s", err.Error())
	}

	return Middleware{
		BotStartTime:  startTime,
		ReceiveChan:   receiveChan,
		SendChan:      sendChan,
		ConsoleChatID: consoleChatID,
		Weather:       w,
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

func readConfig() (int64, error) {
	var consoleChatID int64 = 0
	consoleChatIDStr, err := config.Get(consoleChatIDdir)
	if err != nil {
		return 0, err
	} else {
		log.Printf("Console Chat ID found: %s", consoleChatIDStr)
		consoleChatID, err = strconv.ParseInt(consoleChatIDStr, 10, 64)
		if err != nil {
			return 0, err
		} else {
			log.Printf("Successful parsing of chat ID to int64")
		}
	}

	return consoleChatID, nil
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
		responseTxt := m.prepareResponse(msg.Username, msg.Content)

		if len(responseTxt) > 0 {
			m.SendChan <- Msg{
				ChatID:   msg.ChatID,
				Content:  responseTxt,
				Username: "Ganeungbot",
			}
		}
	}
}

// TODO: come up with a better function name
func (m *Middleware) prepareResponse(username string, txt string) string {
	parseResult := Parse(txt)
	//log.Printf("%+v\n", parseResult)

	var response string
	if len(parseResult.Actions) > 0 {
		response = m.buildResponse(parseResult.Actions[0], username, txt)
	}

	return response
}

func (m *Middleware) buildResponse(action Action, username string, txt string) string {
	var resp = ""
	var err error

	if answerList, ok := Answers[action]; ok {
		resp = util.GetRandomElement(answerList)
	} else if action == ACTION_VERSUS {
		resp = util.GetRandomElement(strings.Split(txt, "vs"))
		log.Printf("%+v\n", util.GetRandomElement(strings.Split(txt, "vs")))
	} else {
		switch action {
		case ACTION_TYPEHELPER:
			resp = typehelper.GetResponse(strings.Split(txt, typehelper.Trigger)[1])
		case ACTION_WEATHER:
			// TODO: send in user's location or user's info so that we can fetch default location per user?
			resp, err = m.Weather.GetResponse(username, txt)
			if err != nil {
				log.Printf(err.Error())
				resp = err.Error()
			}
		//	case ACTION_GAMESTATS:
		//		resp = chatObj.lolStats.GetResponse(txt)

		default:
		}
	}

	//fmt.Printf("%+v\n", resp)
	return resp
}
