package chat

import (
	"bufio"
	"os"

	"github.com/GwangGwang/ganeungbot/pkg/util"
)

// StartConsole instantiates console and returns the channel it'd receive inputs
func StartConsole() chan string {
	scanner := bufio.NewScanner(os.Stdin)
	ch := make(chan string)

	go func() {
		for scanner.Scan() {
			text := scanner.Text()
			if text == "q" {
				close(ch)
				util.Exit()
			} else {
				ch <- text
			}
		}
	}()

	return ch
}
