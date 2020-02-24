package lol

// League of Legends stats via official Riot Games API

import (
	"fmt"
	"log"
)

type (
	LOL struct {
		RiotGamesAPIKey string
		UserInfos []User
		Champions []ChampionInfo
	}

	queryKey struct {
		target targetCategory
		mode gameMode


	}
)


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
	//staticChampionData, err := fetcher.FetchStaticChampionData()
	//if err != nil {
	//	return err
	//}
	//UpsertStaticChampionData(staticChampionData)
	//l.Champions = staticChampionData.Data

	// TODO: game modes

	// 2. summoner
	summoners, err := fetcher.FetchSummoners()
	if err != nil {
		return err
	}

	err = UpsertSummoners(summoners)
	if err != nil {
		return err
	}

	// 3. get matchlists and prepare list of game ids needed to fetch
	// TODO: maybe this map can be map[int64][]string with []string being list of summIds of interest?
	gameIds := make(map[int64]bool)
	for _, summoner := range summoners {
		matchlistRaw, err := fetcher.FetchMatchListRaw(summoner)
		if err != nil {
			return err
		}
		err = UpsertMatchlistRaw(matchlistRaw)
		if err != nil {
			return err
		}

		for _, matchRef := range matchlistRaw.MatchReferencesRaw {
			gameIds[matchRef.GameId] = true
		}
	}

	// 4. get/fetch matches
	matchMap := make(map[int64]MatchRaw)
	log.Printf("fetching total of %d matches\n", len(gameIds))
	for gameId, _ := range gameIds {
		// get match data from either db or url
		matchRaw, err := GetMatchRaw(gameId)
		if err != nil {
			log.Printf("match '%d' not found in db; fetching from url\n", gameId)
			matchRaw, err := fetcher.FetchMatch(gameId)
			if err != nil {
				return err
			}
			err = UpsertMatchRaw(matchRaw)
			if err != nil {
				return err
			}
		}

		// process match data
		matchMap[gameId] = matchRaw
	}

	// 5. process match data per user/ign/queue combination


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

