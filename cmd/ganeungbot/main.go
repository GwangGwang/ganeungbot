package main

import (
	"log"
	"strconv"
	"time"
	"os"

	"github.com/GwangGwang/ganeungbot/pkg/mid"
	"github.com/GwangGwang/ganeungbot/pkg/telegram"
	"github.com/GwangGwang/ganeungbot/pkg/weather"
)

const (
	TelegramApiKey        = "TELEGRAM_API_KEY"
	TelegramConsoleChatId = "TELEGRAM_CONSOLE_CHAT_ID"
	WeatherApiKey         = "WEATHER_API_KEY"
	GeocodingApiKey       = "GEOCODING_API_KEY"
)

var envNames = []string{
	TelegramApiKey,
	TelegramConsoleChatId,
	WeatherApiKey,
	GeocodingApiKey,
}

func main() {
	startTime := time.Now().Unix() - 3600
	log.Printf("Ganeungbot started on %d", startTime)

	envs := make(map[string]string)
	for _, envName := range envNames {
		envs[envName] = os.Getenv(envName)
	}

	// Telegram API
	receiveChan, sendChan, err := telegram.New(envs[TelegramApiKey])
	if err != nil {
		log.Panic(err)
		return
	}

	// Weather API
	w, err := weather.New(envs[WeatherApiKey], envs[GeocodingApiKey])
	if err != nil {
		log.Println(err)
	}

	consoleChatId, err := strconv.ParseInt(envs[TelegramConsoleChatId], 10, 64)
	if err != nil {
		log.Printf("Error while converting consoleChatId to int64: %s", err.Error())
	}

	middleware := mid.Middleware{
		BotStartTime:  startTime,
		ReceiveChan:   receiveChan,
		SendChan:      sendChan,
		ConsoleChatID: consoleChatId,
		Weather:       w,
	}
	middleware.Start()
}
