package lol

// League of Legends stats via official Riot Games API

import (
	"fmt"
	"log"
)

type (
	LOL struct {
		RiotGamesAPIKey  string
		UserInfos        []User
		Champions        []ChampionInfo
		GameIdToMatchRef map[int64][]MatchReference
		AggStats         map[queryKey]aggStats
		BestStats         map[queryKey]bestStats
	}
)

// New initializes and returns a new weather pkg Weather
func New(key string) (LOL, error) {
	log.Println("Initializing lol pkg")

	if len(key) == 0 {
		return LOL{}, fmt.Errorf("riot Games API key not found")
	}

	if err := ensureIndexes(); err != nil {
		return LOL{}, fmt.Errorf("error while ensuring index: %s", err.Error())
	}

	return LOL{
		RiotGamesAPIKey:  key,
		UserInfos:        GetUsers(),
		Champions:        Champions,
		GameIdToMatchRef: make(map[int64][]MatchReference),
		AggStats:         make(map[queryKey]aggStats),
		BestStats:         make(map[queryKey]bestStats),
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
	//l.Champions = staticChampionData.AggStats

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
		matchlist, err := GetMatchlist(summoner)
		if err != nil {
			log.Printf("matchlist for summoner '%s' not found in db; fetching from url\n", summoner.Name)
			matchlist, err := fetcher.FetchMatchList(summoner)
			if err != nil {
				return err
			}
			err = UpsertMatchlist(matchlist)
			if err != nil {
				log.Printf("error while UpsertMatchlist: %s", err.Error())
			}
		}
		log.Printf("got matchlist for summoner '%s'; total of %d matches", summoner.Name, len(matchlist.MatchReferences))

		for _, matchRef := range matchlist.MatchReferences {
			if _, ok := l.GameIdToMatchRef[matchRef.GameId]; !ok {
				// this game id has never been seen before
				l.GameIdToMatchRef[matchRef.GameId] = []MatchReference{}
			}
			matchRef.SummonerName = summoner.Name
			l.GameIdToMatchRef[matchRef.GameId] = append(l.GameIdToMatchRef[matchRef.GameId], matchRef)
		}
	}

	// 4. get/fetch matches
	log.Printf("fetching total of %d matches\n", len(l.GameIdToMatchRef))
	for gameId, matchRefs := range l.GameIdToMatchRef {
		// get match data from either db or url
		match, err := GetMatchRaw(gameId)
		if err != nil || match.GameId == 0 { // entry not found in db
			log.Printf("match '%d' not found in db; fetching from url\n", gameId)
			match, err := fetcher.FetchMatch(gameId)
			if err != nil {
				return err
			}
			err = UpsertMatch(match)
			if err != nil {
				log.Printf("error while UpsertMatch: %s", err.Error())
			}
		}
		log.Printf("got match '%d'", gameId)

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

			// find the index within the data array
			participantIndex := extractParticipantIndex(match, summonerName)
			if participantIndex == -1 {
				log.Printf("summoner '%s' not found in participants of match '%d'", summonerName, gameId)
				continue
			}
			rawStats := match.ParticipantData[participantIndex].Stats
			stats := prepareStats(rawStats)

			queryKeys := createQueryKeys(userName, summonerName, gameMode, champId)

			for _, key := range queryKeys {
				if _, ok := l.AggStats[key]; !ok {
					l.AggStats[key] = aggStats{
						data: make(map[string]float64),
					}
				}
				l.AggStats[key] = updateAvgData(l.AggStats[key], stats)

				key.isBestStats = true
				if _, ok := l.BestStats[key]; !ok {
					l.BestStats[key] = bestStats{
						data: make(map[string]bestData),
					}
				}
				l.BestStats[key] = updateBestData(l.BestStats[key], stats, champIdToChamp(champId))

				// Update overall best data
				// prevent doing this twice by skipping over ingamename ones
				updateOverall := key.target.category == categoryUser

				if updateOverall {
					key.target = target{
						category: categoryAll,
						name: "",
					}
					if _, ok := l.BestStats[key]; !ok {
						// key does not exist; create one
						l.BestStats[key] = bestStats{
							data: make(map[string]bestData),
						}
					}
					var bestOverallData bestStats = l.BestStats[key]
					l.BestStats[key] = updateBestData(bestOverallData, stats, userName)
				}
			}

		} // end process match data for loop

	} // end process all match data for loop

	// 5. pre-generate response texts
	l.prepareResponse()


	// Test code
	for key, val := range l.AggStats {
		fmt.Printf("key: %+v\n", key)
		fmt.Printf("resp:\n%s\n\n", val.resp)
	}
	for key, val := range l.BestStats {
		fmt.Printf("key: %+v\n", key)
		fmt.Printf("resp:\n%s\n\n", val.resp)
	}

	return nil
}

// GetResponse is the main outward facing function to generate response
func (l *LOL) GetResponse(chatID int64, username string, txt string) (string, error) {

	// parse out time/location keywords and process any time offsets
	queryKey, err := l.parse(txt)
	if err != nil {
		return "", err
	}

	fmt.Printf("query key: %+v\n", queryKey)

	if queryKey.isBestStats {
		if data, ok := l.BestStats[queryKey]; ok && data.resp != "" {
			return data.resp, err
		}
	} else {
		if data, ok := l.AggStats[queryKey]; ok && data.resp != "" {
			return data.resp, err
		}
	}

	return "something's wrong; contact admin", nil
}

