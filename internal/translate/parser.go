package translate

import (
	"fmt"
	"regexp"
	"strings"
)

type queryType string

type parseResult struct {
	queryType queryType
	text string
}

const (
	translateText queryType = "translate"
	setLanguage = "setLang"
	getLanguage = "getLang"
	listLanguages = "listLang"
)


type matcher struct {
	regex string
	queryType queryType
}

// order matters here
var matchers = []matcher{
	{ "^번역언어 목록$", listLanguages},
	{"번역언어:", setLanguage},
	{"^번역언어$", getLanguage},
	{"번역:", translateText},
}

const (
	errUnknownQuery = "unknown translate query"
)

func parse(txt string) (parseResult, error) {
	// anonymous helper function for error returning
	handleError := func(err error) (parseResult, error) {
		return parseResult{}, err
	}

	var result parseResult

	for _, matcher := range matchers {
		match, err := regexp.MatchString(matcher.regex, txt)
		if err != nil {
			return handleError(err)
		} else if match {
			result.queryType = matcher.queryType

			if result.queryType == setLanguage || result.queryType == translateText {
				raw := strings.SplitN(txt, matcher.regex, 2)[1]
				result.text = strings.TrimSpace(raw)

				if len(result.text) == 0 {
					return handleError(fmt.Errorf(errUnknownQuery))
				}
			}
			return result, nil
		}

	}

	return handleError(fmt.Errorf(errUnknownQuery))
}
