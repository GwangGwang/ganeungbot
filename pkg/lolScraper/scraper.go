package lolScraper

import (
	"encoding/json"
	"fmt"
	"github.com/GwangGwang/ganeungbot/pkg/lol"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	urlBase = "https://na1.api.riotgames.com/lol/%s?api_key=%s"

	// in seconds
	retryMinor = 1 // 20 requests per second
	retryMajor = 120 // 100 requests per 2 minites (120 seconds)
	retryEscalateCount = 3 // # of times to retry after 1 second before moving on to 20 seconds
)


// Weather is the weather forecast object
type LOLScraper struct {
	RiotGamesAPIKey string
	UserInfos []lol.UserInfo
}

// New initializes and returns a new weather pkg Weather
func New(key string) (LOLScraper, error) {
	log.Println("Initializing lol pkg")

	ins := LOLScraper{}

	if len(key) == 0 {
		return ins, fmt.Errorf("Riot Games API key not found")
	}
	ins.RiotGamesAPIKey = key
	ins.UserInfos = lol.GetUsers()

	return ins, nil
}


func (l *LOLScraper) UpdateMatchList(summonerInfo lol.SummonerInfo) error {
	/*
	  Two problems with matchlists
	  - Only up to 100 returned per query
	  - Query contains the total # of games for the account but this fluctuates depending on beginIndex value
	   (e.g. beginIndex=0 shows 175 games but beginIndex=100 shows 1808 games)
	 */
	summonerName := summonerInfo.Name

	subDomain := fmt.Sprintf("summoner/v4/match/matchlists/by-account/%s", summonerInfo.AccountId)
	subDomain += "?beginIndex=%d" // will be updated iteratively

	// complete url without beginIndex filled yet
	reqUrl := fmt.Sprintf(urlBase, subDomain, l.RiotGamesAPIKey)
	beginIndex := 0
	totalGames := 0

	for {
		log.Printf("retrieving matchlist for summoner '%s'; beginIndex %d", summonerName, beginIndex)
		url := fmt.Sprintf(reqUrl, beginIndex)
		log.Printf("url: %s", url)
		body, err := getWithRetry(url)
		if err != nil {
			return fmt.Errorf("error while retrieving matchlist for summoner '%s': %s", summonerName, err)
		}

		var matchlist Matchlist
		err = json.Unmarshal(body, &matchlist)
		if err != nil {
			return fmt.Errorf("error while unmarshaling matchlist for summoner '%s': %s", summonerName, err)
		}


	}


//	// For some STUPID reason total games count fluctuates depending on beginIndex and need to be checked every loop iteration
//	for {
//		var matchlistFilename string = LOL_DATA_DIR + ingamename + fmt.Sprintf("/matchlist-beginIndex=%d", beginIndex)
//		fmt.Printf("%s - Retrieving matchlist data beginIndex %d for %s\n", LOL_LOG_HEADER, beginIndex, ingamename)
//		payload, err = retrieveData(ingamename, matchlistFilename)
//		if err != nil {
//			log.Print(err)
//		}
//		var matchlistData *MatchList = &MatchList{}
//		err = json.Unmarshal(payload, matchlistData)
//		if err != nil {
//			log.Print(err)
//		}
//
//		fullMatchList = append(fullMatchList, matchlistData.Matches...)
//
//		// Break condition
//		totalGames = matchlistData.TotalGames
//		fmt.Println("%s - Total Games Count: %d", LOL_LOG_HEADER, totalGames)
//		if len(fullMatchList) < totalGames {
//			beginIndex += 100
//		} else {
//			fmt.Printf("%s - Done fetching matchlist data", LOL_LOG_HEADER)
//			break
//		}
//	}



	return nil
}

func (l *LOLScraper) UpdateSummonerInfo() error {
	reqUrl := fmt.Sprintf(urlBase, "summoner/v4/summoners/by-name/%s", l.RiotGamesAPIKey)

	userInfos := lol.GetUsers()

	for _, userInfo := range userInfos {
		for _, summonerName := range userInfo.SummonerNames {
			log.Printf("Retrieving summoner info for summoner '%s'", summonerName)

			summonerUrl := fmt.Sprintf(reqUrl, url.QueryEscape(summonerName))

			body, err := getWithRetry(summonerUrl)
			if err != nil {
				return fmt.Errorf("error while retrieving summoner info for summoner %s: %s", summonerName, err)
			}

			var summonerInfo lol.SummonerInfo
			err = json.Unmarshal(body, &summonerInfo)
			if err != nil {
				return fmt.Errorf("error while unmarshalling summoner info json body for summoner %s: %s", summonerName, err)
			}

			err = UpsertSummonerInfo(summonerInfo)
			if err != nil {
				return fmt.Errorf("error while upserting summoner data for summoner %s: %s", summonerName, err)
			}
		}
	}

	return nil
}

/* STATIC DATA */
func (l *LOLScraper) UpdateStaticChampionData() error {
	url := "http://ddragon.leagueoflegends.com/cdn/9.24.2/data/en_US/champion.json"
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error while retrieving champion data: %s", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var chdata ChampionData
	err = json.Unmarshal(body, &chdata)
	if err != nil {
		return fmt.Errorf("error while unmarshalling champion data json body: %s", err)
	}

	for _, chInfo := range chdata.Data {
		UpsertStaticChampionInfo(chInfo)
	}

	return nil
}

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

