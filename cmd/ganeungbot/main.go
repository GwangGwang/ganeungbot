package main

import (
	"log"
	"strconv"
	"time"

	"github.com/GwangGwang/ganeungbot/pkg/mid"
	"github.com/GwangGwang/ganeungbot/pkg/util"
	"github.com/GwangGwang/ganeungbot/pkg/telegram"
	"github.com/GwangGwang/ganeungbot/pkg/weather"
	"github.com/GwangGwang/ganeungbot/internal/pkg/config"
)

const tokenDir string = "/secrets/telegram"
const consoleChatIDDir string = "/secrets/telegram-consoleChatId"
const weatherAPIKeyDir string = "/secrets/weatherAPIKey"

func main() {
	startTime := time.Now().Unix()
	log.Printf("Ganeungbot started on %d", startTime)

	// Read config
	configMap := config.Get()


	chatIDStr := util.FileRead(consoleChatIDDir)


	// Telegram API
	receiveChan, sendChan, err := telegram.New(configMap[config.TelegramAPIKey])
	if err != nil {
		log.Panic(err)
		return
	}

	// Telegram Console
	consoleChatID, err := strconv.ParseInt(configMap[config.TelegramConsoleChatID], 10, 64)
	if err != nil {
		log.Println("Cannot initialize console: %s", err)
		consoleChatID = 0
	}

	// Weather API
	weather := weather.New(configMap[config.WeatherAPIKey])

	midware := mid.Middleware{
		BotStartTime: startTime,
		ConsoleChatID: consoleChatID,
		ReceiveChan: receiveChan,
		SendChan: sendChan,
		WeatherAPI: weather,
	}

	midware.Start()
}
