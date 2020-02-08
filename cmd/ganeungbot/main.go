package main

import (
	"fmt"

	"log"
	"os"
	"strconv"
	"time"

	"github.com/GwangGwang/ganeungbot/pkg/geocoding"
	"github.com/GwangGwang/ganeungbot/pkg/mid"
	"github.com/GwangGwang/ganeungbot/pkg/telegram"
	"github.com/GwangGwang/ganeungbot/pkg/translate"
	"github.com/GwangGwang/ganeungbot/pkg/weather"
)

const (
	TelegramApiKey        = "TELEGRAM_API_KEY"
	TelegramConsoleChatId = "TELEGRAM_CONSOLE_CHAT_ID"
	WeatherApiKey         = "WEATHER_API_KEY"
	GeocodingApiKey       = "GEOCODING_API_KEY"
	TranslateApiKey       = "TRANSLATE_API_KEY"
)

var envNames = []string{
	TelegramApiKey,
	TelegramConsoleChatId,
	WeatherApiKey,
	GeocodingApiKey,
	TranslateApiKey,
}

func main() {
	startTime := time.Now().Unix() - 3600
	log.Printf("Ganeungbot started on %d", startTime)

	envs := make(map[string]string)
	for _, envName := range envNames {
		envs[envName] = os.Getenv(envName)
	}

	if len(envs[TelegramApiKey]) == 0 {
		panic(fmt.Sprintf("telegram api key not supplied under env '%s'", TelegramApiKey))
	}

	// Telegram API
	receiveChan, sendChan, err := telegram.New(envs[TelegramApiKey])
	if err != nil {
		log.Panic(err)
		return
	}

	// Telegram Console
	consoleChatId, err := strconv.ParseInt(envs[TelegramConsoleChatId], 10, 64)
	if err != nil {
		log.Printf("Error while converting consoleChatId to int64: %s", err.Error())
	}

	// Google Geocoding API
	geocoding, err := geocoding.New(envs[GeocodingApiKey])
	if err != nil {
		log.Printf("Error while initializing geocoding pkg: %s", err.Error())
	}

	// Weather API
	w, err := weather.New(envs[WeatherApiKey], geocoding)
	if err != nil {
		log.Println(err)
	}

	// Translate API
	t, err := translate.New(envs[TranslateApiKey])
	if err != nil {
		log.Println(err)
	}

	middleware := mid.New(startTime, receiveChan, sendChan, consoleChatId, w, t)
	middleware.Start()
}
