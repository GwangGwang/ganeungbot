package lol

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name  string

		// input
		in    string // make sure to append '전적' afterwards
		newChamp championInfo

		// output
		expTargetCategory targetCategory
		expTargetName string
		expGameMode gameMode
		expBestStats bool
		expChampionId int
		expErrStr string // always checked explicitly
	}{
		{name: "all aram stats",
			in:    "전체 아람",
			expTargetCategory: categoryAll,
			expTargetName: "전체",
			expGameMode: gameModeAram,
		},
		{name: "no target",
			in:    "으악 아람",
			expErrStr: errNoTargetFound,
		},
		{name: "invalid game mode",
			in:    "전체 what",
			expErrStr: errNoGameModeKeyword,
		},
		{name: "game mode + best",
			in:    "전체 아람 최고",
			expTargetCategory: categoryAll,
			expTargetName: "전체",
			expGameMode: gameModeAram,
			expBestStats: true,
		},
		{name: "game mode + champion",
			in:    "전체 아람 우왕",
			newChamp: championInfo{id: 1337, name: "blah", matchers: []string{"우왕"}},
			expTargetCategory: categoryAll,
			expTargetName: "전체",
			expGameMode: gameModeAram,
			expChampionId: 1337,
		},
		{name: "game mode + champion + best",
			in:    "전체 아람 우왕 최고",
			newChamp: championInfo{id: 1337, name: "blah", matchers: []string{"우왕"}},
			expTargetCategory: categoryAll,
			expTargetName: "전체",
			expGameMode: gameModeAram,
			expChampionId: 1337,
			expBestStats: true,
		},
		{name: "more than 1 gamemode",
			in:    "전체 아람 노멀",
			expErrStr: errMultipleGameModeKeyword,
		},
		{name: "more than 1 best stats",
			in:    "전체 최고 아람 최고",
			expErrStr: errMultipleBestStatsKeyword,
		},
		{name: "more than 1 champion",
			in:    "전체 우왕 아람 우왕",
			newChamp: championInfo{id: 1337, name: "blah", matchers: []string{"우왕"}},
			expErrStr: errMultipleChampionKeyword,
		},
		{name: "extra unknown words",
			in:    "전체 최고 아람 우왕 바보",
			newChamp: championInfo{id: 1337, name: "blah", matchers: []string{"우왕"}},
			expErrStr: errUnknownKeyword,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// SETUP
			if test.newChamp.id != 0  {
				champions = append(champions, test.newChamp)
			}

			// EXECUTE
			pr, err := parse(test.in + " 전적")

			// VERIFY
			if test.expErrStr == "" {
				if err != nil {
					t.Errorf("expected no error but got '%s'", err.Error())
				} else {
					// parseResult
					if test.expTargetCategory != categoryNone && pr.target.category != test.expTargetCategory {
						t.Errorf("expected target category '%d' but got '%d'", test.expTargetCategory, pr.target.category)
					}
					if test.expTargetName != "" && pr.target.name != test.expTargetName {
						t.Errorf("expected target name '%s' but got '%s'", test.expTargetName, pr.target.name)
					}
					if test.expGameMode != gameModeNone && pr.gameMode != test.expGameMode {
						t.Errorf("expected game mode '%s' but got '%s'", test.expGameMode, pr.gameMode)
					}
					if test.expBestStats != pr.isBestStats {
						t.Errorf("expected isBestStats to be '%t' but got '%t'", test.expBestStats, pr.isBestStats)
					}
					if test.expChampionId != pr.championId {
						t.Errorf("expected championId to be '%d' but got '%d'", test.expChampionId, pr.championId)
					}

				}
			}
			if test.expErrStr != "" && err == nil {
				t.Errorf("expected error '%s' but got nil", getError(test.expErrStr))
			}
		})
	}
}


func TestParseTarget(t *testing.T) {
	tests := []struct {
		name  string

		// input
		in    string
		newUser map[string]userinfo

		// output
		targetCategory targetCategory
		targetName string
		out string
	}{
		{name: "all",
			in:    "전체 아람",
			targetCategory: categoryAll,
			targetName: "전체",
			out: "아람",
		},
		{name: "user",
			in:    "광승 아람",
			targetCategory: categoryUser,
			targetName: "광승",
			out: "아람",
		},
		{name: "ign",
			in:    "KwangKwang 아람",
			targetCategory: categoryIgn,
			targetName: "KwangKwang",
			out: "아람",
		},
		{name: "username found in non-first word",
			in:    "아람 KwangKwang",
			targetCategory: categoryNone,
			targetName: "",
			out: "아람 KwangKwang",
		},
		{name: "more than 2 words",
			in:    "KwangKwang 아람 루시안",
			targetCategory: categoryIgn,
			targetName: "KwangKwang",
			out: "아람 루시안",
		},
		{name: "multi word username",
			in:    "im gosu 아람 루시안",
			newUser: map[string]userinfo{
				"im gosu": {igns: []string{"what"}},
			},
			targetCategory: categoryUser,
			targetName: "im gosu",
			out: "아람 루시안",
		},
		{name: "multi word ign",
			in:    "what the 아람 루시안",
			newUser: map[string]userinfo{
				"im gosu": {igns: []string{"what the"}},
			},
			targetCategory: categoryIgn,
			targetName: "what the",
			out: "아람 루시안",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// SETUP
			if test.newUser != nil {
				for name, info := range test.newUser {
					usermap[name] = info
				}
			}

			// EXECUTE
			actualTarget, actualWords := parseTarget(strings.Split(test.in, " "))

			// VERIFY
			// target
			if test.targetCategory != actualTarget.category {
				t.Errorf("expected target category '%d' but found '%d'", test.targetCategory, actualTarget.category)
			}
			if test.targetName != actualTarget.name {
				t.Errorf("expected target name '%s' but found '%s'", test.targetName, actualTarget.name)
			}

			// output words
			expWords := strings.Split(test.out, " ")
			if len(expWords) != len(actualWords) {
				t.Errorf("expected # of words to be %d but found %d", len(expWords), len(actualWords))
			} else {
				for i, expWord := range expWords {
					if expWord != actualWords[i] {
						t.Errorf("expected word index %d to be %s but found %s", i, expWord, actualWords[i])
					}

				}
			}


		})
	}
}

