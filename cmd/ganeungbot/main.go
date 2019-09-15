package main

import (
	"log"
	"strconv"
	"time"

	"github.com/GwangGwang/ganeungbot/pkg/mid"
	"github.com/GwangGwang/ganeungbot/pkg/util"
)

const tokenDir string = "/secrets/telegram"
const consoleChatIDDir string = "/secrets/telegram-consoleChatId"

func main() {
	startTime := time.Now().Unix()
	log.Printf("Ganeungbot started on %d", startTime)

	// Read config
	// TODO: move to internal
	token := util.FileRead(tokenDir)
	chatIDStr := util.FileRead(consoleChatIDDir)
	consoleChatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	util.Check(err)

	receiveChan, sendChan, err := telegram.New(token)
	if err != nil {
		log.Panic(err)
		return
	}
	consoleChan := mid.StartConsole()

	mid.Init(receiveChan, sendChan, consoleChan)

}
