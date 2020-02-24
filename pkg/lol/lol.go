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
		GameIdToMatchRef map[int64][]MatchReferenceRaw
		Data map[queryKey]aggData
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
		Champions: Champions,
		GameIdToMatchRef: make(map[int64][]MatchReferenceRaw),
		Data: make(map[queryKey]aggData),
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

	// every game contains one or more users in this chat group
	// keep track so that we don't have to parse a match more than once for users
	for _, summoner := range summoners {
		matchlistRaw, err := GetMatchlistRaw(summoner)
		if err != nil {
			log.Printf("matchlist for summoner '%d' not found in db; fetching from url\n", summoner.Name)
			matchlistRaw, err := fetcher.FetchMatchListRaw(summoner)
			if err != nil {
				return err
			}
			err = UpsertMatchlistRaw(matchlistRaw)
			if err != nil {
				log.Printf("error while UpsertMatchlistRaw: %s", err.Error())
			}
		}

		for _, matchRef := range matchlistRaw.MatchReferencesRaw {
			if _, ok := l.GameIdToMatchRef[matchRef.GameId]; !ok {
				// this game id has never been seen before
				l.GameIdToMatchRef[matchRef.GameId] = []MatchReferenceRaw{}
			}
			matchRef.SummonerName = summoner.Name
			l.GameIdToMatchRef[matchRef.GameId] = append(l.GameIdToMatchRef[matchRef.GameId], matchRef)
		}
	}

	// 4. get/fetch matches
	log.Printf("fetching total of %d matches\n", len(l.GameIdToMatchRef))
	for gameId, matchRefs := range l.GameIdToMatchRef {
		// get match data from either db or url
		matchRaw, err := GetMatchRaw(gameId)
		if err != nil || matchRaw.GameId == 0 { // entry not found in db
			log.Printf("match '%d' not found in db; fetching from url\n", gameId)
			matchRaw, err := fetcher.FetchMatch(gameId)
			if err != nil {
				return err
			}
			err = UpsertMatchRaw(matchRaw)
			if err != nil {
				log.Printf("error while UpsertMatchRaw: %s", err.Error())
			}
		}

		/* process match data
		 1. identify which users participated in this particular match (via saved matchRefs list)
		 2. create queryKeys for indices in aggregated data that need to be updated
		   note that each match per user will contain user / game mode / champ, meaning the updates are for:
		  - avg/summoner/champ
		  - avg/summoner/all
		  - avg/user/champ
		  - avg/user/all
		  - best/summoner/champ
		  - best/summoner/all
		  - best/user/champ
		  - best/user/all
		 3. update (or create) dataset for each queryKey, aggregating for avg and replacing with highest for best
		 */
		for _, matchRef := range matchRefs {
			gameMode := getQueueFromId(matchRef.QueueId)
			if gameMode == gameModeNone {
				continue
			}
			champId := matchRef.ChampionId
			summonerName := matchRef.SummonerName
			userName := getUserFromSummoner(l.UserInfos, summonerName)

			queryKeys := createQueryKeys(userName, summonerName, gameMode, champId)

		}


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

