package weather

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	// error msgs - Syntax
	errUnexpectedWord    = "no more words expected but found"
	errUnexpectedTimeKeyword = "unexpected time keyword"
	errMoreThanOneLocation  = "location words showed up at more than one place in the query"

	// error msgs - Semantic
	err24HrTimeAboveCeiling = "hour greater than 23 not allowed"
	err12HrTimeAboveCeiling = "only up to 11 allowed for am/pm time"
	errUnknownTimeToken     = "unknown time token type"
	errPastTime             = "request is for time in the past"

	// time constants
	secondsInHour int64 = 3600
	secondsInDay int64 = 86400
	defaultAM int64 = 9  // 9am
	defaultPM int64 = 16 // 4pm
)

/* parse out time keywords and location phrase to be used in geocoding/forecast apis
 @param txt string   : the string to parse
 @return parseResult : object encapsulating time query type, time offset, and the location
 @return error       : error; should be nil for valid query
*/
func parse(txt string) (parseResult, error) {
	// anonymous helper function for error returning
	handleError := func(err error) (parseResult, error) {
		return parseResult{}, err
	}

	// 1. validate msg structure and parse out time keywords and location
	timeTokens, location, err := parseSyntax(txt)
	if err != nil {
		return handleError(err)
	}

	// default to querying for current forecast if no time-related keywords are found
	var queryType tokenType = now
	var timeOffset int64 = 0
	if len(timeTokens) > 0 {
		// 2. semantically validate time keywords and compute which time to query for
		curHour := int64(time.Now().Hour())
		queryType, timeOffset, err = parseSemantic(timeTokens, curHour)
		if err != nil {
			return handleError(err)
		}
	}

	return parseResult{
		TimeInfo: timeinfo{
			Category: queryType,
			Offset: timeOffset,
		},
		Location: location,
	}, nil
}

/* parse out time keywords and location phrase, while also validating their syntax
 @param txt string : the string to parse
 @return []token   : list of validated time keywords encapsulated in token object
 @return string    : location phrase to query geocoding
 @return error     : error; should be nil for valid query
 */
func parseSyntax(txt string) ([]token, string, error) {
	// anonymous helper function for error returning
	handleError := func(word string, errString string) ([]token, string, error) {
		return []token{}, "", getParseError(word, errString)
	}

	// split request string into words and take off the last part (날씨)
	words := strings.Split(txt, " ")
	words = words[0 : len(words)-1] // Remove the useless 날씨

	/* at least in the scope of this function, location words are any words that aren't time keywords
	location words has to be grouped together, i.e. queries like 'toronto 6pm ON' makes no sense
	the two booleans denote:
	 locationStart: true if there has ever been a location word; any additional location word will be appended to existing collection of location words
	 locationEnd: becomes true when locationStart was true and a time keyword showed up, meaning we shouldn't see any more non-time keywords
	*/
	var locationWords = make([]string, 0)
	var locationStart, locationEnd = false, false

	/* Main parsing logic
	1. Search for time keyword, from the list of known time words
	2. Upon finding a time keyword, check that this time keyword is part of the expected token type list (currentTokenTypes);
	 if so, add the time keyword to validated time keyword list and update next possible token list

	 While doing so...
	- The first time finding a non-time keyword, set locationStart to true and start collection location words
	- Upon finding a time keyword again, set the collected location words as the location to query for
	*/
	validatedList := make([]token, 0)
	currentTokenTypes := allTokenTypes
	for _, word := range words {
		// Not expecting time nor location word
		if len(currentTokenTypes) == 0 && locationEnd {
			return handleError(word, errUnexpectedWord)
		}

		// Loop through all currently available token types and check if this word belongs to any
		// match bool is used as a flag outside of loop and also to skip any more unnecessary looping after match is found
		match := false
		for _, tokenType := range allTokenTypes {
			tokenInfo := timeTokenMap[tokenType]
			for _, tokenGroup := range tokenInfo.TokenGroups {
				for _, keyword := range tokenGroup.Keywords {
					match, _ = regexp.MatchString(fmt.Sprintf("^(?i)%s$", keyword), word)
					if match {
						// check if this tokenType belongs to expected types, and error out otherwise
						validToken := false
						for _, expTokenType := range currentTokenTypes {
							if tokenType == expTokenType {
								validToken = true
								break
							}
						}
						if !validToken {
							return handleError(word, errUnexpectedTimeKeyword)
						}

						// token is part of the expected list of token types; add to validated list and update next list of tokens to expect
						validatedList = append(validatedList, token{
							word: word,
							tokenType: tokenType,
							value: tokenGroup.Value,
						})
						currentTokenTypes = tokenInfo.NextTokens
						break
					}
				}
				if match {
					break
				}
			}
			if match {
				break
			}
		}
		if match {
			// this word is a time keyword; check if we were in midst of searching for location words, and mark end it if so
			if locationStart && !locationEnd {
				locationEnd = true
			}
		} else {
			// this word is a non-time keyword, meaning it has to be location word
			if !locationStart {
				// marks start of location words
				locationStart = true
				locationWords = append(locationWords, word)
			} else if locationStart && !locationEnd {
				// In the midst of grabbing location words
				locationWords = append(locationWords, word)
			} else if locationEnd {
				// already finished gathering location but another non-time word showed up
				return handleError(word, errMoreThanOneLocation)
			}
		}
	}

	// Join back the location words
	finalLocation := strings.Join(locationWords, " ")

	return validatedList, finalLocation, nil
}

func getParseError(word string, errString string) error {
	return fmt.Errorf("Error at '%s': %s", word, errString)
}

/* parse and verify semantics of time keywords retrieved from syntax parse
 assumption: time tokens are well-formed syntactically
 @param tokens []token : the list of time keyword tokens to verify semantics
 @return tokenType     : type of query when requesting the forecast api, with regards to time
 @return int64         : time offset to be applied on top of 0am today for time query
 @return error         : error; should be nil for valid query
*/
func parseSemantic(tokens []token, curHour int64) (tokenType, int64, error) {
	// anonymous helper function for error returning
	handleError := func(word string, errString string) (tokenType, int64, error) {
		return "", 0, getParseError(word, errString)
	}

	// query type will be based on whatever the last time keyword's type
	var queryType tokenType = now
	if len(tokens) > 0 {
		queryType = tokens[len(tokens)-1].tokenType
		if queryType == ampm || queryType == hourAmpm {
			queryType = hour
		}
	}

	var offset int64 = 0
	var ampmVal int64 = -1
	isToday := true

	hourRegex := regexp.MustCompile(`\d+`)
	for i, token := range tokens {
		word := token.word
		switch token.tokenType {
		case now:
			// no-op
		case day:
			if token.value != 0 {
				isToday = false
			}
			offset += token.value * secondsInDay
		case ampm:
			ampmVal = token.value

			// asked for TODAY am time when currently is pm
			if isToday && curHour >= 12 && token.value == 0 {
				return handleError(word, errPastTime)
			}

			// additional handling if this is the last token
			if i == len(tokens)-1 {
				var queryHour int64 = 0
				// defaulting to 9am and 4pm for am/pm
				if token.value == 0 {
					queryHour = defaultAM
				} else if token.value == 1 {
					queryHour = defaultPM
				}
				// if the defaulted time is of the past for today, query for whatever +1 hour from present
				if (isToday && curHour >= queryHour) {
					queryHour = curHour + 1
				}

				offset += queryHour * secondsInHour
			} else {
				offset += token.value * secondsInDay / 2
			}
		case hour:
			hourNum, err := strconv.ParseInt(hourRegex.FindString(word), 10, 64)
			if err != nil {
				return handleError(word, err.Error())
			}

			if hourNum >= 24 || hourNum < 0 {
				return handleError(word, err24HrTimeAboveCeiling)
			} else if ampmVal != -1 && hourNum >= 12 {
				// query is in the format of "오전 x시" or "오후 x시"; check if value is in [1,12]
				return handleError(word, err12HrTimeAboveCeiling)
			}

			queryHour := hourNum
			if isToday {
				/* Special handling for hour forecasts without am/pm specified
				1. For explicit am/pm hour query, the time must be in the future
				2. For hour query for today, try to round up to nearest explicit am/pm and error otherwise
				ex) now is 4am, request for at 5 --> use 5am as intended
				    now is 4pm, request for at 5 --> round up to 5pm
				    now is 4pm, request for at 3 --> invalid (next 3 is 3am tomorrow)
				*/
				if ampmVal == -1 {
					if curHour < queryHour {
						// good
					} else if curHour > 12 && queryHour < 12 && curHour < queryHour+12 {
						queryHour += 12
					} else {
						return handleError(word, errPastTime)
					}
				} else if ampmVal == 1 {
					if curHour >= queryHour + 12 {
						return handleError(word, errPastTime)
					}
				}
			}
			offset += queryHour * secondsInHour
		case hourAmpm:
			hourNum, err := strconv.ParseInt(hourRegex.FindString(word), 10, 64)
			if err != nil {
				return handleError(word, err.Error())
			}
			if hourNum >= 12 {
				return handleError(word, err12HrTimeAboveCeiling)
			} else {
				queryHour := token.value*12 + hourNum

				if isToday && queryHour <= curHour {
					return handleError(word, errPastTime)
				}
				// add additional 12 hours if pm
				offset += queryHour * secondsInHour
			}
		}
	}

	return queryType, offset, nil
}

func getSavedLocation(loc string) string {
	var savedLoc string
	if len(loc) > 0 {
		locationMatched := false
		for _, locationParseGroup := range locationParseGroups {
			for _, keyword := range locationParseGroup.Keywords {
				if matched, _ := regexp.MatchString(fmt.Sprintf("(?i)^%s$", keyword), loc); matched {
					savedLoc = locationParseGroup.Value
					locationMatched = true
					break
				}
			}
			if locationMatched {
				break
			}
		}
	}

	// TODO: retrieve default location per user

	return savedLoc
}

