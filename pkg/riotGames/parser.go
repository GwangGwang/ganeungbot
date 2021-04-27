package riotGames

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	// error msgs - Syntax
	errNoTargetFound = "no stats target found from query; needs to be one of '전체' / username / in game name"
	errQueryTooShort = "query too short; needs to be minimum 3 words, including target / game mode"
	errNoGameModeKeyword = "no game mode keyword found (normal, ranked, aram, etc.)"
	errMultipleBestStatsKeyword = "more than one keyword for best stats"
	errMultipleGameModeKeyword = "more than one keyword for game mode"
	errMultipleChampionKeyword = "more than one keyword for champion name"
	errUnknownKeyword = "unknown keyword"
)

/*
  Query Criteria

 * target audience: whose stats?
   - ign
   - person (i.e. multiple igns)
   - everyone (for jangin stats)
 * game type
   - normal
   - ranked
   - rift (normal + ranked)
   - aram
   - bot?
 * (optional) best stats
 * (optional) champion

 [Target] [GameMode] (champion) (best) 전적/stats
 ex) 유저 노멀 루시안 전적
     전체 아람 전적
     전체 랭크 최고 전적

need collections:
 - user to igns relations
 - champion info
 */

/* parse out time keywords and location phrase to be used in geocoding/forecast apis
 @param txt string   : the string to parse
 @return queryKey : object encapsulating query results
 @return error       : error; should be nil for valid query
*/
func (l *LOL) parse(txt string) (queryKey, error) {
	// anonymous helper function for error returning
	handleError := func(errMsg string) (queryKey, error) {
		return queryKey{}, getError(errMsg)
	}

	result := queryKey{}

	// split request string into words
	words := strings.Split(txt, " ")

	// query needs to be minimum of 3 words (target gamemode stats)
	if len(words) < 3 {
		return handleError(errQueryTooShort)
	} else {
		words = words[:len(words)-1] // take out the last token (trigger token)
	}

	// 1. parse target
	target, words := l.parseSubject(words)
	if target.category == categoryNone {
		return handleError(errNoTargetFound)
	} else {
		result.target = target
	}

	// 2. parse game mode / (champion) / (best)
	for _, word := range words {
		if word == "최고" {
			// prevent duplicate keyword
			if result.isBestStats {
				return handleError(errMultipleBestStatsKeyword)
			}
			result.isBestStats = true
		} else if mode := extractGameMode(word); mode != gameModeNone {
			if result.gameMode != gameModeNone {
				return handleError(errMultipleGameModeKeyword)
			}
			result.gameMode = mode
		} else if champ := l.extractChampion(word); champ.Id != "" {
			if result.championId != "" {
				return handleError(errMultipleChampionKeyword)
			}
			result.championId = champ.Id
		} else {
			return handleError(errUnknownKeyword + " " + word)
		}
	}

	return result, nil
}

func extractGameMode(word string) gameMode {
	for _, queue := range queues {
		for _, matcher := range queue.matchers {
			if word == matcher {
				return queue.gameMode
				break
			}
		}
	}

	return gameModeNone
}

func (l *LOL) extractChampion(word string) ChampionInfo {
	for _, chInfo := range l.Champions {
		for _, matcher := range chInfo.Matchers {
			match, _ := regexp.MatchString(matcher, word)
			if match {
				return chInfo
			}
		}
	}

	return ChampionInfo{}
}

// returns target found and words left after parsing out the target
func (l *LOL) parseSubject(words []string) (target, []string) {
	// TODO: should add defaulting to querying user if no target found

	result := target{}

	// 1. parse target
	if words[0] == "전체" {
		result = target{
			category: categoryAll,
			name: "",
		}
		words = words[1:]
	} else {
		for _, userinfo := range l.UserInfos {
			user := userinfo.HumanName
			userwords := strings.Split(user, " ")
			matched := true
			for i, userword := range userwords {
				if words[i] != userword {
					matched = false
					break
				}
			}
			if matched {
				result = target{
					category: categoryUser,
					name: userinfo.UserName,
				}
				words = words[len(userwords):]
				break
			} else {
				ignMatched := false
				for _, ign := range userinfo.SummonerNames {
					ignwords := strings.Split(ign, " ")
					matched := true
					for i, ignword := range ignwords {
						if words[i] != ignword {
							matched = false
							break
						}
					}
					if matched {
						result = target{
							category: categoryIgn,
							name:     ign,
						}
						words = words[len(ignwords):]
						ignMatched = true
						break
					}
				}
				if ignMatched {
					break
				}
			}
		}
	}

	return result, words
}

func getError(errMsg string) error {
	return fmt.Errorf("LOL stats parsing error: %s", errMsg)
}

