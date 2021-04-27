package riotGames

import (
	"fmt"
	"strings"
)

var respLines map[string]string = map[string]string{
	"Kills":               "KDA: %.1f/",
	"Deaths":              "%.1f/",
	"Assists":             "%.1f ",
	"KDA":                 "(avg %.2f)\n",
	"DmgPerLife":          "Damage per Life: %.2f\n",
	"TotalDmgChamp":       "- Total Damage Dealt: %.2f\n",
	"DmgTakenPerLife":     "Damage taken per Life: %.2f\n",
	"TotalDamageTaken":    "- Total Damage Taken: %.2f\n",
	"DamageSelfMitigated": "- Damage Self-Mitigated: %.2f\n",
	"TotalHeal":           "Total Healing Done: %.2f\n",
	"CS":                  "CS: %.2f\n",
	"Vision Score":        "Vision Score: %.2f\n",
}

func (l *LOL) prepareResponse() {
	for key, AggStats := range l.AggStats {
		var stats map[string]float64 = AggStats.data
		var respLineAvg map[string]string = respLines
		var resp string

		if stats["GamesPlayed"] == 0 {
			resp += "No games played\n"
		} else {
			var numGames float64 = stats["GamesPlayed"]
			stats["Winrate"] = stats["Win"] * float64(100) // will be divided by game count in for loop
			respLineAvg["Winrate"] = "Winrate: %.2f%% " + fmt.Sprintf("(%.0f Games)\n", numGames)

			for _, statsKey := range []string{
				"Winrate",
				"Kills", "Deaths", "Assists", "KDA",
				"DmgPerLife", "TotalDmgChamp",
				"DmgTakenPerLife", "TotalDamageTaken", "DamageSelfMitigated",
				"TotalHeal",
				"CS",
				"Vision Score",
			} {
				if key.gameMode == gameModeAram && statsKey == "Vision Score" {
					continue
				}

				if statsKey == "KDA" {
					stats[statsKey] = computeKDA(stats["Kills"], stats["Deaths"], stats["Assists"])
				} else {
					stats[statsKey] /= numGames
				}
				resp += fmt.Sprintf(respLineAvg[statsKey], stats[statsKey])
			}

		}
		var loldata = aggStats{
			data: stats,
			resp: resp,
		}
		l.AggStats[key] = loldata
	}

	// 4-2. best data per user (e.g. 광승 루시안 노멀 최고 전적, 광승 노멀 최고 전적)
	for key, bestData := range l.BestStats {
		var numGames float64 = l.AggStats[key].data["GamesPlayed"]
		stats := bestData.data
		var respLineBest map[string]string = respLines
		var resp string

		if numGames == 0 {
			resp += "No games played\n"
		} else if key.target.category != categoryAll {
			// redefine some of the lines to print
			respLineBest["KDA"] = "KDA: %.2f " + fmt.Sprintf("(%.0f/%.0f/%.0f)\n", bestData.kda.Kills, bestData.kda.Deaths, bestData.kda.Assists)
			respLineBest["Kills"] = "- Kills: %.0f\n"
			delete(respLineBest, "Deaths")
			respLineBest["Assists"] = "- Assists: %.0f\n"

			for _, statsKey := range []string{
				"KDA",
				"Kills", "Assists",
				"DmgPerLife", "TotalDmgChamp",
				"DmgTakenPerLife",
				"TotalDamageTaken",
				"DamageSelfMitigated",
				"TotalHeal",
				"CS",
				"Vision Score",
			} {
				if key.gameMode == gameModeAram && statsKey == "Vision Score" {
					continue
				}
				if line, ok := respLineBest[statsKey]; ok {
					var respline string = fmt.Sprintf(line, stats[statsKey].value)
					if key.championId != "0" {
						// Show the best champ/user
						respline = strings.Replace(respline, "\n", generateBestDataSuffix(stats[statsKey].candidates), 1)
					}
					resp += respline
				}

			}

		}
		// Add Games Played in front for all champ stats
		//		if key.championId != 0 {
		//			resp = fmt.Sprintf("Games Played: %.0f (%s)", numGames, generateBestDataSuffix(stats["GamesPlayed"].candidates))
		//		}

		var lolbestdata = bestStats{
			data: stats,
			resp: resp,
		}
		l.BestStats[key] = lolbestdata
	}


}
