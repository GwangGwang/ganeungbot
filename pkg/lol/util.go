package lol

import (
	"fmt"
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
