package translate

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name  string
		valid bool

		// input
		in    string

		// output
		queryType queryType
		text string
		errstr string
	}{
		{name: "listLanguages",
			in:    "번역언어 목록",
			valid: true,
			queryType: listLanguages,
		},
		{name: "listLanguages fail",
			in:    "번역언어 목록이당",
		},
		{name: "setLanguage",
			in:    "번역언어:ko",
			valid: true,
			queryType: setLanguage,
			text: "ko",
		},
		{name: "setLanguage with whitespace",
			in:    "번역언어:  ko",
			valid: true,
			queryType: setLanguage,
			text: "ko",
		},
		{name: "setLanguage - edge case of repeating query phrase",
			in:    "번역언어: 번역언어:   번역언어:",
			valid: true,
			queryType: setLanguage,
			text: "번역언어:   번역언어:",
		},
		{name: "setLanguage empty",
			in:    "번역언어:",
		},
		{name: "getLanguage",
			in:    "번역언어",
			valid: true,
			queryType: getLanguage,
		},
		{name: "getLanguage fail",
			in:    "번역언어달라",
		},
		{name: "translate",
			in:    "번역: hi",
			valid: true,
			queryType: translateText,
			text: "hi",
		},
		{name: "translate - whitespaces",
			in:    "번역: hi i'm boy",
			valid: true,
			queryType: translateText,
			text: "hi i'm boy",
		},
		{name: "translate - empty",
			in:    "번역:",
		},
		{name: "translate - edge case of repeating query phrase",
			in:    "번역:번역: hi 번역:i'm boy   ",
			valid: true,
			queryType: translateText,
			text: "번역: hi 번역:i'm boy",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// EXECUTE
			actual, err := parse(test.in)

			// VERIFY
			if test.valid {
				assert.NoError(t, err)
			} else if !test.valid && assert.Error(t, err) {
				assert.EqualError(t, err, errUnknownQuery)
			}

			// even for error cases these should assert true since empty values
			assert.Equal(t, test.queryType, actual.queryType)
			assert.Equal(t, test.text, actual.text)

		})
	}
}
