package riotGames

import (
	"fmt"
	"github.com/fatih/structs"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func createQueryKeys(username string, summonerName string, mode gameMode, champId int) []queryKey {
	return []queryKey{
		{
			// summoner - mode - champ
			target: target{
				category: categoryIgn,
				name: summonerName,
			},
			gameMode: mode,
			isBestStats: false,
			championId: fmt.Sprintf("%d", champId),
		},
		{
			// summoner - mode - all
			target: target{
				category: categoryIgn,
				name: summonerName,
			},
			gameMode: mode,
			isBestStats: false,
			championId: "",
		},
		{
			// user - mode - champ
			target: target{
				category: categoryUser,
				name: username,
			},
			gameMode: mode,
			isBestStats: false,
			championId: fmt.Sprintf("%d", champId),
		},
		{
			// user - mode - all
			target: target{
				category: categoryUser,
				name: username,
			},
			gameMode: mode,
			isBestStats: false,
			championId: "",
		},
	}
}

func extractParticipantIndex(match Match, summonerName string) int {
	for i, pi := range match.ParticipantIdentities {
		if pi.Data.SummonerName == summonerName {
			return i
		}
	}

	return -1
}

// Convert ParticipantStats object to string->float64 map and add more data such as kda, cs, etc.
func prepareStats(stats ParticipantStats) map[string]float64 {
	var mapStats map[string]float64 = convertToMap(stats)

	mapStats["KDA"] = computeKDA(mapStats["Kills"], mapStats["Deaths"], mapStats["Assists"])
	mapStats["CS"] = computeCS(mapStats["TotalMinionsKilled"], mapStats["NeutralKilled"])
	mapStats["DmgPerLife"] = computeDmgPerLife(mapStats["TotalDamageDealtToChampions"], mapStats["Deaths"])
	mapStats["DmgTakenPerLife"] = computeDmgTakenPerLife(mapStats["TotalDamageTaken"], mapStats["DamageSelfMitigated"], mapStats["Deaths"])

	//fmt.Printf("*** prepareMapStats *** \n%+v", mapStats)
	return mapStats
}

func convertToMap(stats ParticipantStats) map[string]float64 {
	var resultMap = make(map[string]float64)
	var mapStats = structs.Map(stats)

	for field, val := range mapStats {
		var newVal float64
		if f64val, ok := val.(float64); ok {
			newVal = f64val
		} else if intVal, ok := val.(int); ok {
			newVal = float64(intVal)
		} else if int64Val, ok := val.(int64); ok {
			newVal = float64(int64Val)
		} else if boolVal, ok := val.(bool); ok {
			if boolVal {
				newVal = 1
			}
		}
		resultMap[field] = newVal
	}

	return resultMap
}

func computeKDA(kills float64, deaths float64, assists float64) float64 {
	if deaths == 0 {
		deaths = 1 / 1.2 // No death means multiply K+A by 1.2
	}
	return (kills + assists) / deaths
}

func computeCS(minions float64, neutrals float64) float64 {
	return minions + neutrals
}

func computeDmgPerLife(dmg float64, deaths float64) float64 {
	return dmg / (deaths + 1)
}

func computeDmgTakenPerLife(dmgtaken float64, mitig float64, deaths float64) float64 {
	return (dmgtaken + mitig) / (deaths + 1)
}

func updateAvgData(origdata aggStats, stats map[string]float64) aggStats {
	var gamedata aggStats = origdata
	gamedata.data["GamesPlayed"] += 1

	for _, field := range []string{
		"GamesPlayed",
		"Win",
		"Kills",
		"Deaths",
		"Assists",
		"TotalDamageDealtToChampions",
		"DmgPerLife",
		"TotalHeal",
		"TotalDamageTaken",
		"DamageSelfMitigated",
		"DmgTakenPerLife",
		"TotalMinionsKilled",
		"NeutralKilled",
		"CS",
		"VisionScore",
	} {
		gamedata.data[field] += stats[field]
	}
	//fmt.Printf("*** updateAvgData *** \n%+v", gamedata.data)

	return gamedata
}

func updateBestData(origdata bestStats, stats map[string]float64, bestCandidate string) bestStats {
	var data = origdata.data

	for _, field := range []string{
		"GamesPlayed",
		"Win",
		"Kills",
		"Assists",
		"KDA",
		"TotalDamageDealtToChampions",
		"DmgPerLife",
		"TotalHeal",
		"TotalDamageTaken",
		"DamageSelfMitigated",
		"DmgTakenPerLife",
		"CS",
		"VisionScore",
	} {
		var newVal float64 = stats[field]

		if data == nil {
			data = make(map[string]bestData)
		}

		if oldBest, ok := data[field]; ok {
			// previous best value exists; compare and replace (or append)
			if newVal == oldBest.value {
				oldBest.candidates[bestCandidate] = struct{}{}
			} else if newVal > oldBest.value {
				oldBest.candidates = make(map[string]struct{})
				oldBest.candidates[bestCandidate] = struct{}{}
				oldBest.value = newVal

				if field == "KDA" {
					origdata.kda = updateBestKDA(stats)
				}
			}
			data[field] = oldBest
		} else {
			// no best data so far
			//var kda KDA = KDA{ Kills:
			cand := make(map[string]struct{})
			cand[bestCandidate] = struct{}{}
			data[field] = bestData{value: newVal, candidates: cand}
			if field == "KDA" {
				origdata.kda = updateBestKDA(stats)
			}
		}
	}

	origdata.data = data
	return origdata
}

func updateBestKDA(stats map[string]float64) KDA {
	return KDA{
		Kills:   stats["Kills"],
		Deaths:  stats["Deaths"],
		Assists: stats["Assists"],
	}
}


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

func generateBestDataSuffix(cands map[string]struct{}) string {
	var resp string = " (%s)\n"
	var candString string // "cand1,cand2,cand3"
	var count int
	for k := range cands {
		if candString != "" {
			candString += ","
		}
		candString += k

		count++

		if count >= 5 {
			break
		}
	}

	if len(cands) > 5 {
		candString += ",etc.)"
	}

	return fmt.Sprintf(resp, candString)
}
