package mid

import (
	"regexp"
)

// ParseResult is the result of parsing, such as tokens, etc.
type ParseResult struct {
	Actions []Action
}

// Parse reads and deduces data from msg
func Parse(txt string) ParseResult {
	actions := scanForActions(txt)

	parseResult := ParseResult{
		Actions: actions,
	}

	return parseResult
}

func scanForActions(txt string) []Action {
	var actions []Action
	for _, cmd := range Commands {

		for _, keyword := range cmd.Keywords {
			// (?i) - case insensitive
			match, _ := regexp.MatchString("(?i)"+keyword, txt)

			if match {
				actions = append(actions, cmd.Behavior)
			}

		}

	}

	return actions
}
