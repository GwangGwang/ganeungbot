package lol

import (
	"encoding/json"
	"log"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	urlBase = "https://na1.api.riotgames.com/lol/%s?api_key=%s"

	// in seconds
	retryMinor = 1 // 20 requests per second
	retryMajor = 120 // 100 requests per 2 minitues (120 seconds)
	retryEscalateCount = 3 // # of times to retry after 1 second before moving on to 20 seconds
)

func (l *LOL) UpdateSummonerInfo() error {
	reqUrl := fmt.Sprintf(urlBase, "summoner/v4/summoners/by-name/%s",l.RiotGamesAPIKey)

	userInfos := GetUsers()

	for _, userinfo := range userInfos {
		for _, summonerName := range userinfo.SummonerNames {
			log.Printf("Retrieving summoner info for summoner '%s'", summonerName)

			summonerUrl := fmt.Sprintf(reqUrl, url.QueryEscape(summonerName))

			resp, err := getWithRetry(summonerUrl)
			if err != nil {
				return fmt.Errorf("error while retrieving summoner info for summoner %s: %s", summonerName, err)
			}

			body, err := ioutil.ReadAll(resp.Body)

			var summonerInfo SummonerInfo
			err = json.Unmarshal(body, &summonerInfo)
			if err != nil {
				return fmt.Errorf("error while unmarshalling summoner info json body for summoner %s: %s", summonerName, err)
			}

			err = UpsertSummonerInfo(summonerInfo)
			if err != nil {
				return fmt.Errorf("error while upserting summoner data for summoner %s: %s", summonerName, err)
			}
			resp.Body.Close()
		}
	}

	return nil
}


/* STATIC DATA */
func (l *LOL) UpdateStaticChampionData() error {
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
func getWithRetry(reqUrl string) (*http.Response, error) {
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
			waitRetry(retryAmount)
		} else {
			return resp, err
		}
	}
}

func waitRetry(sec int) {
	log.Printf("rate limited; retrying in %d seconds", sec)
	time.Sleep(time.Duration(sec) * time.Second)
}

