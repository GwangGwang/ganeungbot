package lol

// League of Legends stats via official Riot Games API

import (
	"fmt"
	"log"
)

type LOL struct {
	RiotGamesAPIKey string
	UserInfoMap map[int64]string
	UserInfos []User
	Champions map[string]ChampionInfo
}

// New initializes and returns a new weather pkg Weather
func New(key string) (LOL, error) {
	log.Println("Initializing lol pkg")

	if len(key) == 0 {
		return LOL{}, fmt.Errorf("Riot Games API key not found")
	}

	return LOL{
		RiotGamesAPIKey: key,
		UserInfos: GetUsers(),
	}, nil
}

func (l *LOL) Update() error {
	fetcher, err := NewFetcher(l.RiotGamesAPIKey)
	if err != nil {
		return err
	}

	// 1. static data
	staticChampionData, err := fetcher.FetchStaticChampionData()
	if err != nil {
		return err
	}
	//UpsertStaticChampionData(staticChampionData)
	l.Champions = staticChampionData.Data

	// TODO: game modes

	// 2. summoner
	summonersData, err := fetcher.FetchSummoners()
	if err != nil {
		return err
	}

	err = UpsertSummoners(summonersData)
	if err != nil {
		return err
	}

	// 3. matchlists




	return nil
}

// GetResponse is the main outward facing function to generate response
func (l *LOL) GetResponse(chatID int64, username string, txt string) (string, error) {

	// parse out time/location keywords and process any time offsets
	_, err := l.parse(txt)
	if err != nil {
		return "", err
	}

	return "", nil
}

