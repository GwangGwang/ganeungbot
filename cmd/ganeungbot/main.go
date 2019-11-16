package main

import (
	"log"
	"time"

	"github.com/GwangGwang/ganeungbot/pkg/mid"
	"github.com/GwangGwang/ganeungbot/pkg/telegram"
	"github.com/GwangGwang/ganeungbot/pkg/weather"
)

const tokenDir string = "/secrets/telegram"
const consoleChatIDDir string = "/secrets/telegram-consoleChatId"
const weatherAPIKeyDir string = "/secrets/weatherAPIKey"

func main() {
	startTime := time.Now().Unix() - 3600
	log.Printf("Ganeungbot started on %d", startTime)

	// Telegram API
	receiveChan, sendChan, err := telegram.New()
	if err != nil {
		log.Panic(err)
		return
	}

	// Weather API
	w, err := weather.New()
	if err != nil {
		log.Println(err)
	}

	midware := mid.New(startTime, receiveChan, sendChan, w)
	midware.Start()
}
