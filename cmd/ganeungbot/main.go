package main

import (
	"log"
	"strconv"
	"time"

	"github.com/GwangGwang/ganeungbot/pkg/chat"
	"github.com/GwangGwang/ganeungbot/pkg/telegram"
	"github.com/GwangGwang/ganeungbot/pkg/util"
)

const tokenDir string = "/secrets/telegram"
const consoleChatIDDir string = "/secrets/telegram-consoleChatId"

func main() {
	startTime := time.Now().Unix()
	log.Printf("Ganeungbot started on %d", startTime)

	// Read config
	token := util.FileRead(tokenDir)
	chatIDStr := util.FileRead(consoleChatIDDir)
	consoleChatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	util.Check(err)

	// Uses generic chat interface so that supporting another msg service is trivial
	consoleChan := chat.StartConsole()
	var chatObj chat.Chat
	chatObj = &telegram.Telegram{
		Token:          token,
		ConsoleChatID:  consoleChatID,
		ConsoleChannel: consoleChan,
	}
	chatObj.Start()
}
