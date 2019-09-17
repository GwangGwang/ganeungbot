package mid

import (
	"bufio"
	"os"

	"github.com/GwangGwang/ganeungbot/pkg/util"
)

// startConsole instantiates console and returns the channel it'd receive inputs
func startConsole(sendChan chan Msg, consoleChatID int64) {
	scanner := bufio.NewScanner(os.Stdin)

	go func() {
		for scanner.Scan() {
			text := scanner.Text()
			if text == "q" {
				close(sendChan)
				util.Exit()
			} else {
				msg := Msg{
					ChatID:   consoleChatID,
					Username: "GaneungBot",
					Content:  text,
				}
				sendChan <- msg
			}
		}
	}()
}
