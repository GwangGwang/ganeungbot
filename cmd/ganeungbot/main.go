package main

import (
	"log"
	"time"

	"github.com/GwangGwang/ganeungbot/pkg/console"
	"github.com/GwangGwang/ganeungbot/pkg/telegram"
)

func main() {
	startTime := time.Now().Unix()
	log.Printf("Ganeungbot started on %d", startTime)

	consoleChan := console.Start()

	telegram.InitBot(consoleChan)
}
