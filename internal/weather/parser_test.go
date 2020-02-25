package weather

import (
	"regexp"
	"testing"
)

func TestParseSyntax(t *testing.T) {
	tests := []struct {
		name  string
		valid bool

		// input
		in    string

		// output
		tokens   []token
		location string
		err      error
	}{
		{name: "Empty",
			in:    "날씨",
			valid: true,
		},
		{name: "now",
			in:       "현재 날씨",
			valid:    true,
			tokens:   []token{ {
				word:      "현재",
				tokenType: now,
			} },
		},
		{name: "another time keyword after 'now' keywords",
			in:  "현재 9시 날씨",
			err: getParseError("9시", errUnexpectedTimeKeyword),
		},
		{name: "day - today",
			in:    "오늘 날씨",
			valid: true,
			tokens:   []token{ {
				word:      "오늘",
				tokenType: day,
				value:     0,
			} },
		},
		{name: "day - tomorrow",
			in:     "내일 날씨",
			valid:  true,
			tokens:   []token{ {
				word:      "내일",
				tokenType: day,
				value:     1,
			} },
		},
		{name: "day - day after tomorrow",
			in:     "모레 날씨",
			valid:  true,
			tokens:   []token{ {
				word:      "모레",
				tokenType: day,
				value:     2,
			} },
		},
		{name: "day - 2 day after tomorrow",
			in:     "글피 날씨",
			valid:  true,
			tokens:   []token{ {
				word:      "글피",
				tokenType: day,
				value:     3,
			} },
		},
		{name: "day keyword followed by another",
			in:  "오늘 내일 날씨",
			err: getParseError("내일", errUnexpectedTimeKeyword),
		},
		{name: "day + am",
			in:     "내일 오전 날씨",
			valid:  true,
			tokens:   []token{ {
				word:      "내일",
				tokenType: day,
				value:     1,
			}, {
				word:      "오전",
				tokenType: ampm,
				value:     0,
			} },
		},
		{name: "day + pm + hour",
			in:     "내일 오후 2시 날씨",
			valid:  true,
			tokens:   []token{ {
				word:      "내일",
				tokenType: day,
				value:     1,
			}, {
				word:      "오후",
				tokenType: ampm,
				value:     1,
			}, {
				word:      "2시",
				tokenType: hour,
			},
		} },
		{name: "day + hourAmpm",
			in:     "내일 3am 날씨",
			valid:  true,
			tokens:   []token{ {
				word:      "내일",
				tokenType: day,
				value:     1,
			}, {
				word:      "3am",
				tokenType: hourAmpm,
				value:     0,
			} },
		},
		{name: "am",
			in:     "오전 날씨",
			valid:  true,
			tokens:   []token{ {
				word:      "오전",
				tokenType: ampm,
				value:     0,
			} },
		},
		{name: "pm",
			in:     "오후 날씨",
			valid:  true,
			tokens:   []token{ {
				word:      "오후",
				tokenType: ampm,
				value:     1,
			} },
		},
		{name: "pm + hour",
			in:     "오후 10시 날씨",
			valid:  true,
			tokens:   []token{ {
				word:      "오후",
				tokenType: ampm,
				value:     1,
			}, {
				word:      "10시",
				tokenType: hour,
			} },
		},
		{name: "hour",
			in:     "22시 날씨",
			valid:  true,
			tokens:   []token{ {
				word:      "22시",
				tokenType: hour,
			} },
		},
		{name: "hour + hour",
			in:  "9시 12시 날씨",
			err: getParseError("12시", errUnexpectedTimeKeyword),
		},
		{name: "hourAm",
			in:     "8am 날씨",
			valid:  true,
			tokens:   []token{ {
				word:      "8am",
				tokenType: hourAmpm,
				value:     0,
			} },
		},
		{name: "hourPm",
			in:     "8pm 날씨",
			valid:  true,
			tokens:   []token{ {
				word:      "8pm",
				tokenType: hourAmpm,
				value:     1,
			} },
		},
		{name: "hourPm + hourPm",
			in:  "8pm 5am 날씨",
			err: getParseError("5am", errUnexpectedTimeKeyword),
		},
		{name: "location before and after time keyword",
			in:  "토론토 오늘 밴쿠버 날씨",
			err: getParseError("밴쿠버", errMoreThanOneLocation),
		},
		{name: "Single non-time keyword input",
			in:       "Toronto 날씨",
			valid:    true,
			location: "Toronto",
		},
		{name: "Multiple non-time keyword input",
			in:       "Toronto ON 날씨",
			valid:    true,
			location: "Toronto ON",
		},
		{name: "location + time keywords",
			in:       "asdf tomorrow 날씨",
			valid:    true,
			location: "asdf",
			tokens:   []token{ {
				word:      "tomorrow",
				tokenType: day,
				value:     1,
			} },
		},
		{name: "time + location keywords",
			in:       "모레 토론토 날씨",
			valid:    true,
			location: "토론토",
			tokens:   []token{ {
				word:      "모레",
				tokenType: day,
				value:     2,
			} },
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// EXECUTE
			actualTokens, actualLocation, actualErr := parseSyntax(test.in)

			// VERIFY
			if test.valid {
				if actualErr != nil {
					t.Errorf("Expected no error but error found: %s", actualErr.Error())
				}

				expectedTokens := make([]token, 0)
				expectedLocation := ""
				if test.tokens != nil {
					expectedTokens= test.tokens
				}
				if test.location != "" {
					expectedLocation = test.location
				}

				if len(expectedTokens) != len(actualTokens) {
					t.Errorf("Expected %d tokens parsed but only found %d;\nexpected:\n%+v\n,actual:\n%+v", len(expectedTokens), len(actualTokens), expectedTokens, actualTokens)
				} else {
					for i, expToken := range expectedTokens {
						if expToken != actualTokens[i] {
							t.Errorf("mismatch in token index %d:\nexpected:\n%+v\n,actual:\n%+v", i, expToken, actualTokens[i])
						}
					}
				}
				if expectedLocation != actualLocation {
					t.Errorf("expected location '%s' but found '%s'", expectedLocation, actualLocation)
				}
			} else if !test.valid {
				if actualErr == nil {
					t.Errorf("Expected error but no error found")
				} else if matched, _ := regexp.MatchString(test.err.Error(), actualErr.Error()); !matched {
					t.Errorf("Expected error '%s' but got another: %s", test.err.Error(), actualErr.Error())
				}
			}

		})
	}
}

func TestParseSemantic(t *testing.T) {
	tests := []struct {
		name  string
		valid bool

		// input
		in []token
		curHour int64

		// output
		tokenType tokenType
		offset    int64
		err       error
	}{
		{name: "Empty",
			valid: true,
			tokenType: now,
		},
		{name: "now",
			in: []token{ {
				word:      "현재",
				tokenType: now,
			} },
			valid:    true,
			tokenType: now,
		},
		{name: "day - today",
			in: []token{ {
				tokenType: day,
				value:     0,
			} },
			valid: true,
			tokenType: day,
		},
		{name: "day - tomorrow",
			in: []token{ {
				tokenType: day,
				value:     1,
			} },
			valid: true,
			tokenType: day,
			offset: secondsInDay,
		},
		{name: "day + am",
			in: []token{ {
				word:      "내일",
				tokenType: day,
				value:     1,
			}, {
				word:      "오전",
				tokenType: ampm,
				value:     0,
			} },
			valid: true,
			tokenType: hour,
			offset: secondsInDay + defaultAM * secondsInHour, // defaults to 9am tmr
		},
		{name: "day + pm + hour",
			in: []token{{
				word:      "내일",
				tokenType: day,
				value:     1,
			}, {
				word:      "오후",
				tokenType: ampm,
				value:     1,
			}, {
				word:      "2시",
				tokenType: hour,
			}},
			valid:     true,
			tokenType: hour,
			offset: secondsInDay + (12+2)*secondsInHour,
		},
		{name: "day + hourAmpm",
			in: []token{ {
				word:      "내일",
				tokenType: day,
				value:     1,
			}, {
				word:      "3am",
				tokenType: hourAmpm,
				value:     0,
			} },
			valid:  true,
			tokenType: hour,
			offset: secondsInDay + 3*secondsInHour,
		},
		{name: "am, current time is after default",
			valid:  true,
			in:   []token{ {
				word:      "오전",
				tokenType: ampm,
				value:     0,
			} },
			curHour: defaultAM + 1,
			tokenType: hour,
			offset: (defaultAM + 2) * secondsInHour,
		},
		{name: "am, current time is before default",
			valid:  true,
			in:   []token{ {
				word:      "오전",
				tokenType: ampm,
				value:     0,
			} },
			curHour: defaultAM - 2,
			tokenType: hour,
			offset: defaultAM * secondsInHour,
		},
		{name: "am, current time is on default",
			valid:  true,
			in:   []token{ {
				word:      "오전",
				tokenType: ampm,
				value:     0,
			} },
			curHour: defaultAM,
			tokenType: hour,
			offset: (defaultAM + 1) * secondsInHour,
		},
		{name: "pm, current time is after default",
			valid:  true,
			in:   []token{ {
				word:      "오후",
				tokenType: ampm,
				value:     1,
			} },
			curHour: defaultPM + 12 + 2,
			tokenType: hour,
			offset: (defaultPM + 12 + 3) * secondsInHour,
		},
		{name: "pm, current time is before default",
			valid:  true,
			in:   []token{ {
				word:      "오후",
				tokenType: ampm,
				value:     1,
			} },
			curHour: defaultPM - 2,
			tokenType: hour,
			offset: defaultPM * secondsInHour,
		},
		{name: "pm, current time is on default",
			valid:  true,
			in:   []token{ {
				word:      "오후",
				tokenType: ampm,
				value:     1,
			} },
			curHour: defaultPM,
			tokenType: hour,
			offset: (defaultPM + 1) * secondsInHour,
		},
		{name: "pm + hour, before current time",
			in:   []token{ {
				word:      "오후",
				tokenType: ampm,
				value:     1,
			}, {
				word:      "10시",
				tokenType: hour,
			} },
			curHour: 23,
			err: getParseError("10시", errPastTime),
		},
		{name: "pm + hour, after current time",
			valid:  true,
			in:   []token{ {
				word:      "오후",
				tokenType: ampm,
				value:     1,
			}, {
				word:      "10시",
				tokenType: hour,
			} },
			curHour: 19,
			tokenType: hour,
			offset: 22 * secondsInHour,
		},
		{name: "pm + hour, on current time",
			in:   []token{ {
				word:      "오후",
				tokenType: ampm,
				value:     1,
			}, {
				word:      "7시",
				tokenType: hour,
			} },
			curHour: 19,
			err: getParseError("7시", errPastTime),
		},
		{name: "pm + hour > 12",
			in:   []token{ {
				word:      "오후",
				tokenType: ampm,
				value:     1,
			}, {
				word:      "12시",
				tokenType: hour,
			} },
			curHour: 11,
			err: getParseError("12시", err12HrTimeAboveCeiling),
		},
		{name: "am + hour > 12",
			in:   []token{ {
				word:      "오전",
				tokenType: ampm,
				value:     0,
			}, {
				word:      "13시",
				tokenType: hour,
			} },
			curHour: 11,
			err: getParseError("13시", err12HrTimeAboveCeiling),
		},
		{name: "hour > 12, after current time",
			valid:  true,
			in:   []token{ {
				word:      "22시",
				tokenType: hour,
			} },
			curHour: 21,
			tokenType: hour,
			offset: 22 * secondsInHour,
		},
		{name: "hour > 12, before current time",
			in:   []token{ {
				word:      "22시",
				tokenType: hour,
			} },
			curHour: 23,
			err: getParseError("22시", errPastTime),
		},
		{name: "hour > 12, on current time",
			in:   []token{ {
				word:      "22시",
				tokenType: hour,
			} },
			curHour: 22,
			err: getParseError("22시", errPastTime),
		},
		{name: "hour < 12 but if switched to pm, after current time",
			valid: true,
			in:   []token{ {
				word:      "2시",
				tokenType: hour,
			} },
			curHour: 13,
			tokenType: hour,
			offset: 14 * secondsInHour,
		},
		{name: "hour < 12 and if switched to pm, still before current time",
			in:   []token{ {
				word:      "2시",
				tokenType: hour,
			} },
			curHour: 15,
			err: getParseError("2시", errPastTime),
		},
		{name: "hour < 12 and if switched to pm, on current time",
			in:   []token{ {
				word:      "2시",
				tokenType: hour,
			} },
			curHour: 14,
			err: getParseError("2시", errPastTime),
		},
		{name: "hourAm, after current time",
			valid:  true,
			in:   []token{ {
				word:      "8am",
				tokenType: hourAmpm,
				value:     0,
			} },
			curHour: 7,
			tokenType: hour,
			offset: 8 * secondsInHour,
		},
		{name: "hourPm, before current time",
			in:   []token{ {
				word:      "8pm",
				tokenType: hourAmpm,
				value:     1,
			} },
			curHour: 21,
			err: getParseError("8pm", errPastTime),
		},
		{name: "hourPm, on current time",
			in:   []token{ {
				word:      "8pm",
				tokenType: hourAmpm,
				value:     1,
			} },
			curHour: 20,
			err: getParseError("8pm", errPastTime),
		},
		{name: "hourPm, after current time",
			valid:  true,
			in:   []token{ {
				word:      "8pm",
				tokenType: hourAmpm,
				value:     1,
			} },
			curHour: 19,
			tokenType: hour,
			offset: 20 * secondsInHour,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// EXECUTE
			if test.in == nil {
				test.in = make([]token, 0)
			}
			actualTokenType, actualOffset, actualErr := parseSemantic(test.in, test.curHour)

			// VERIFY
			if test.valid {
				if actualErr != nil {
					t.Errorf("Expected no error but error found: %s", actualErr.Error())
				}

				var expectedTokenType tokenType = now
				if test.tokenType != none {
					expectedTokenType= test.tokenType
				}

				if expectedTokenType != actualTokenType {
					t.Errorf("expected token type output '%s' but found '%s'", expectedTokenType, actualTokenType)
				}
				if test.offset != actualOffset {
					t.Errorf("expected offset output '%d' but found '%d'", test.offset, actualOffset)
				}
			} else if !test.valid {
				if actualErr == nil {
					t.Errorf("Expected error but no error found")
				} else if matched, _ := regexp.MatchString(test.err.Error(), actualErr.Error()); !matched {
					t.Errorf("Expected error\n'%s'\nbut got another:\n'%s'", test.err.Error(), actualErr.Error())
				}
			}

		})
	}
}
