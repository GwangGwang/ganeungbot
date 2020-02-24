package lol

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	urlBase = "https://na1.api.riotgames.com/lol/%sapi_key=%s"

	// in seconds
	retryMinor = 1 // 20 requests per second
	retryMajor = 120 // 100 requests per 2 minites (120 seconds)
	retryEscalateCount = 3 // # of times to retry after 1 second before moving on to 20 seconds
)

type Fetcher struct {
	RiotGamesAPIKey string
	UserInfos []User
}

// New initializes and returns a new weather pkg Weather
func NewFetcher(key string) (Fetcher, error) {
	log.Println("Initializing lol pkg")

	ins := Fetcher{}

	if len(key) == 0 {
		return ins, fmt.Errorf("Riot Games API key not found")
	}
	ins.RiotGamesAPIKey = key
	ins.UserInfos = GetUsers()

	return ins, nil
}

func (l *Fetcher) FetchSummoners() ([]Summoner, error) {
	reqUrl := fmt.Sprintf(urlBase, "summoner/v4/summoners/by-name/%s?", l.RiotGamesAPIKey)

	userInfos := GetUsers()

	var summoners []Summoner
	for _, userInfo := range userInfos {
		for _, summonerName := range userInfo.SummonerNames {
			log.Printf("Retrieving summoner info for summoner '%s'", summonerName)

			summonerUrl := fmt.Sprintf(reqUrl, url.QueryEscape(summonerName))

			log.Printf("url: %s", summonerUrl)
			body, err := getWithRetry(summonerUrl)
			if err != nil {
				return []Summoner{}, fmt.Errorf("error while retrieving summoner info for summoner %s: %s", summonerName, err)
			}

			var summoner Summoner
			err = json.Unmarshal(body, &summoner)
			if err != nil {
				return []Summoner{}, fmt.Errorf("error while unmarshalling summoner info json body for summoner %s: %s", summonerName, err)
			}

			summoners = append(summoners, summoner)
		}
	}

	return summoners, nil
}

func (l *Fetcher) FetchMatchListRaw(summonerInfo Summoner) (MatchlistRaw, error) {
	/*
	  Two problems with matchlists
	  - Only up to 100 returned per query
	  - Query contains the total # of games for the account but this fluctuates depending on beginIndex value
	   (e.g. beginIndex=0 shows 175 games but beginIndex=100 shows 1808 games)
	 */
	summonerName := summonerInfo.Name

	subDomain := fmt.Sprintf("match/v4/matchlists/by-account/%s", summonerInfo.AccountId)
	subDomain += "?beginIndex=%d&" // will be updated iteratively

	// complete url without beginIndex filled yet
	reqUrl := fmt.Sprintf(urlBase, subDomain, l.RiotGamesAPIKey)
	beginIndex := 0
	totalGames := 0
	var allMatchReferences []MatchReferenceRaw

	for {
		log.Printf("retrieving matchlist for summoner '%s'; beginIndex %d", summonerName, beginIndex)
		url := fmt.Sprintf(reqUrl, beginIndex)
		log.Printf("url: %s", url)
		body, err := getWithRetry(url)
		if err != nil {
			return MatchlistRaw{}, fmt.Errorf("error while retrieving matchlist for summoner '%s': %s", summonerName, err)
		}

		var matchlistRaw MatchlistRaw
		err = json.Unmarshal(body, &matchlistRaw)
		if err != nil {
			return MatchlistRaw{}, fmt.Errorf("error while unmarshaling matchlist for summoner '%s': %s", summonerName, err)
		}

		allMatchReferences = append(allMatchReferences, matchlistRaw.MatchReferencesRaw...)

		// break condition
		totalGames = matchlistRaw.TotalGames
		log.Printf("total games count according to this batch: %d\n", totalGames)
		if len(allMatchReferences) < totalGames {
			// should retrieve next (max) 100 games
			beginIndex += 100
		} else {
			// this is the last batch
			log.Printf("done fetching total of %d matchlist data for summoner '%s'\n", len(allMatchReferences), summonerName)
			break
		}
	}

	return MatchlistRaw{
		SummonerName: summonerInfo.Name,
		MatchReferencesRaw: allMatchReferences,
	}, nil
}

func (l *Fetcher) FetchMatch(gameId int64) (MatchRaw, error) {
	subDomain := fmt.Sprintf("match/v4/matches/%d?", gameId)
	url := fmt.Sprintf(urlBase, subDomain, l.RiotGamesAPIKey)

	log.Printf("retrieving match id '%d'\n", gameId)
	log.Printf("url: %s\n", url)
	body, err := getWithRetry(url)
	if err != nil {
		return MatchRaw{}, fmt.Errorf("error while retrieving match id '%d'", gameId)
	}

	var matchRaw MatchRaw
	err = json.Unmarshal(body, &matchRaw)
	if err != nil {
		return MatchRaw{}, fmt.Errorf("error while unmarshaling match id '%d'", gameId)
	}

	return matchRaw, nil
}

/* STATIC DATA */
func (l *Fetcher) FetchStaticChampionData() (ChampionDataRaw, error) {
	url := "http://ddragon.leagueoflegends.com/cdn/9.24.2/data/en_US/champion.json"
	resp, err := http.Get(url)
	if err != nil {
		return ChampionDataRaw{}, fmt.Errorf("error while retrieving champion data: %s", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var chdata ChampionDataRaw
	err = json.Unmarshal(body, &chdata)
	if err != nil {
		return ChampionDataRaw{}, fmt.Errorf("error while unmarshalling champion data json body: %s", err)
	}

	return chdata, nil
}

/*
func (l *Fetcher) FetchStaticQueueData() (QueueDataRaw, error) {
	url := "http://static.developer.riotgames.com/docs/lol/queues.json"
	resp, err := http.Get(url)
	if err != nil {
		return QueueDataRaw{}, fmt.Errorf("error while retrieving queue data: %s", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var queueData QueueDataRaw
	err = json.Unmarshal(body, &chdata)
	if err != nil {
		return QueueDataRaw{}, fmt.Errorf("error while unmarshalling queue data json body: %s", err)
	}

	return queueData, nil
}
*/

/* HELPERS */

// Riot API is heavily rate-limited; wait for rate limit to be lifted and retry
func getWithRetry(reqUrl string) ([]byte, error) {
	retryCount := 0
	for {
		resp, err := http.Get(reqUrl)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			retryCount += 1
			retryAmount := retryMinor
			if retryCount >= retryEscalateCount {
				retryCount = 0
				retryAmount = retryMajor
			}
			waitRetry(retryAmount) // TODO: use retryAmount instead of hardcode
		} else if resp.StatusCode == http.StatusOK {
			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			return body, err
		} else {
			return nil, fmt.Errorf("non-200 response status: %s", resp.Status)
		}
	}
}

func waitRetry(sec int) {
	log.Printf("rate limited; retrying in %d seconds", sec)
	time.Sleep(time.Duration(sec) * time.Second)
}

