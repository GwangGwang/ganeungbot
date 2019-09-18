package util

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Exit terminates the program
func Exit() {
	// TODO: save before exit?
	os.Exit(3)
}

// Check reviews error object
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

// FileRead reads file into string
func FileRead(filename string) string {
	data, err := ioutil.ReadFile(filename)
	Check(err)
	return strings.TrimSpace(string(data))
}

// PrintChatLog prints received chat update
func PrintChatLog(chatID int64, msgID int, username string, text string) {
	log.Printf("ChatID:%d | MsgID: %d | %s | %s", chatID, msgID, username, text)
}

// GetRandomElement returns a random element from the given list
func GetRandomElement(arr []string) string {
	rand.Seed(time.Now().Unix()) // TODO: should be called only once
	return arr[rand.Intn(len(arr))]
}
