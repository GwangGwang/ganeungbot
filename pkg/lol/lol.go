package lol

// League of Legends stats via official Riot Games API

import (
	"fmt"
	"log"
)

// Weather is the weather forecast object
type LOL struct {
	RiotGamesAPIKey string
	UserInfos []UserInfo
}

// New initializes and returns a new weather pkg Weather
func New(key string) (LOL, error) {
	log.Println("Initializing lol pkg")

	ins := LOL{}

	if len(key) == 0 {
		return ins, fmt.Errorf("Riot Games API key not found")
	}
	ins.RiotGamesAPIKey = key
	ins.UserInfos = GetUsers()

	return ins, nil
}

// GetResponse is the main outward facing function to generate response
func (l *LOL) GetResponse(username string, txt string) (string, error) {
	// parse out time/location keywords and process any time offsets
	_, err := l.parse(txt)
	if err != nil {
		return "", err
	}

	return "", nil
}

