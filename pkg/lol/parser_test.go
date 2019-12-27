package lol

import (
	"strings"
	"testing"
)

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

