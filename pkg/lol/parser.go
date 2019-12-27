package lol

import (
	"fmt"
	"strings"
)

const (
	// error msgs - Syntax
	errNoTargetFound = "no stats target found from query; needs to be one of '전체' / username / in game name"
	errQueryTooShort = "query too short; needs to be minimum 3 words, including target / game mode"
	errNoGameModeKeyword = "no game mode keyword found (normal, ranked, aram, etc.)"
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

need collections:
 - user to igns relations
 - champion info
 */

/* parse out time keywords and location phrase to be used in geocoding/forecast apis
 @param txt string   : the string to parse
 @return parseResult : object encapsulating query results
 @return error       : error; should be nil for valid query
*/
func parse(txt string) (parseResult, error) {
	// anonymous helper function for error returning
	handleError := func(errMsg string) (parseResult, error) {
		return parseResult{}, fmt.Errorf("LOL stats parsing error: %s", errMsg)
	}

	result := parseResult{}

	// split request string into words
	words := strings.Split(txt, " ")

	// query needs to be minimum of 3 words (target gamemode stats)
	if len(words) < 3 {
		return handleError(errQueryTooShort)
	} else {
		words = words[:len(words)-1] // take out the last token (trigger token)
	}

	// 1. parse target
	target, words := parseTarget(words)
	if target.category == categoryNone {
		return handleError(errNoTargetFound)
	}






		//	for _, queue := range queues {

		//		for _, matcher := range queue.matchers {
		//			if keyword == matcher {
		//				result.gameMode = queue.gameMode
		//				break
		//			}
		//		}
		//	}

	if result.gameMode == none {
		return handleError(errNoGameModeKeyword)
	}

	return result, nil

}

func parseTarget(words []string) (target, []string) {

	result := target{}

	// 1. parse target
	if words[0] == "전체" {
		result = target{
			category: categoryAll,
			name: "전체",
		}
		words = words[1:]
	} else {
		for user, userinfo := range usermap {
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
					name: user,
				}
				words = words[len(userwords):]
				break
			} else {
				ignMatched := false
				for _, ign := range userinfo.igns {
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

