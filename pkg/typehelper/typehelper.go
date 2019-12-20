package typehelper

// Typehelper, mapping English keys to Korean and vice versa

import (
	"fmt"
	"strings"
)

// Char to trigger typehelping
const Trigger = "`"

var ENG_KEY = "rRseEfaqQtTdwWczxvgkoiOjpuPhynbml"
var KOR_KEY = []string{"ㄱ","ㄲ","ㄴ","ㄷ","ㄸ","ㄹ","ㅁ","ㅂ","ㅃ","ㅅ","ㅆ","ㅇ","ㅈ","ㅉ","ㅊ","ㅋ","ㅌ","ㅍ","ㅎ","ㅏ","ㅐ","ㅑ","ㅒ","ㅓ","ㅔ","ㅕ","ㅖ","ㅗ","ㅛ","ㅜ","ㅠ","ㅡ","ㅣ"}
var CHO_DATA = []string{"ㄱ","ㄲ","ㄴ","ㄷ","ㄸ","ㄹ","ㅁ","ㅂ","ㅃ","ㅅ","ㅆ","ㅇ","ㅈ","ㅉ","ㅊ","ㅋ","ㅌ","ㅍ","ㅎ"}
var JUNG_DATA = []string{"ㅏ","ㅐ","ㅑ","ㅒ","ㅓ","ㅔ","ㅕ","ㅖ","ㅗ","ㅘ","ㅙ","ㅚ","ㅛ","ㅜ","ㅝ","ㅞ","ㅟ","ㅠ","ㅡ","ㅢ","ㅣ"}
var JONG_DATA = []string{"ㄱ","ㄲ","ㄳ","ㄴ","ㄵ","ㄶ","ㄷ","ㄹ","ㄺ","ㄻ","ㄼ","ㄽ","ㄾ","ㄿ","ㅀ","ㅁ","ㅂ","ㅄ","ㅅ","ㅆ","ㅇ","ㅈ","ㅊ","ㅋ","ㅌ","ㅍ","ㅎ"}

func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1    //not found.
}


func GetResponse(username string, txt string) string {
	// parse out time/location keywords and process any time offsets
	resp := engTypeToKor(txt)
	return fmt.Sprintf("%s: %s", username, resp)
}


func engTypeToKor(txt string) string {
	// 초성, 중성, 종성 개수
	nCho, nJung, nJong := -1, -1, -1
	var res string

	for _, ch := range txt {
		p := strings.IndexRune(ENG_KEY, ch)

		if p == -1 {
			// not eng
			if nCho != -1 {
				if nJung != -1 { // 초+중+(종)성

					res += makeHangul(nCho, nJung, nJong)
				} else { // 초성만
					res += CHO_DATA[nCho]
				}
			} else {
				if nJung != -1 { // 중성만
					res += JUNG_DATA[nJung]
				} else if nJong != -1 { // 복자음
					res += JONG_DATA[nJong]
				}
			}
			nCho, nJung, nJong = -1, -1, -1
			res += string(ch)
		} else if p < 19 {
			if nJung != -1 {
				if nCho == -1 { // 중성만 입력됨, 초성으로
					res += JUNG_DATA[nJung]
					nJung = -1
					nCho = indexOf(KOR_KEY[p], CHO_DATA)
				} else {             // 종성이다
					if nJong == -1 { // 종성 입력 중
						nJong = indexOf(KOR_KEY[p], JONG_DATA)
						if nJong == -1 { // 종성이 아니라 초성이다
							res += makeHangul(nCho, nJung, nJong)
							nCho = indexOf(KOR_KEY[p], CHO_DATA)
							nJung = -1
						}
					} else if nJong == 0 && p == 9 { // ㄳ
						nJong = 2
					} else if nJong == 3 && p == 12 { // ㄵ
						nJong = 4
					} else if nJong == 3 && p == 18 { // ㄶ
						nJong = 5
					} else if nJong == 7 && p == 0 { // ㄺ
						nJong = 8
					} else if nJong == 7 && p == 6 { // ㄻ
						nJong = 9
					} else if nJong == 7 && p == 7 { // ㄼ
						nJong = 10
					} else if nJong == 7 && p == 9 { // ㄽ
						nJong = 11
					} else if nJong == 7 && p == 16 { // ㄾ
						nJong = 12
					} else if nJong == 7 && p == 17 { // ㄿ
						nJong = 13
					} else if nJong == 7 && p == 18 { // ㅀ
						nJong = 14
					} else if nJong == 16 && p == 9 { // ㅄ
						nJong = 17
					} else {						// 종성 입력 끝, 초성으로
						res += makeHangul(nCho, nJung, nJong)
						nCho = indexOf(KOR_KEY[p], CHO_DATA)
						nJung = -1
						nJong = -1
					}
				}
			} else {            // 초성 또는 (단/복)자음이다
				if nCho == -1 { // 초성 입력 시작
					if nJong != -1 { // 복자음 후 초성
						res += JONG_DATA[nJong]
						nJong = -1
					}
					nCho = indexOf(KOR_KEY[p], CHO_DATA)
				} else if nCho == 0 && p == 9 { // ㄳ
					nCho = -1
					nJong = 2
				} else if nCho == 2 && p == 12 { // ㄵ
					nCho = -1
					nJong = 4
				} else if nCho == 2 && p == 18 { // ㄶ
					nCho = -1
					nJong = 5
				} else if nCho == 5 && p == 0 { // ㄺ
					nCho = -1
					nJong = 8
				} else if nCho == 5 && p == 6 { // ㄻ
					nCho = -1
					nJong = 9
				} else if nCho == 5 && p == 7 { // ㄼ
					nCho = -1
					nJong = 10
				} else if nCho == 5 && p == 9 { // ㄽ
					nCho = -1
					nJong = 11
				} else if nCho == 5 && p == 16 { // ㄾ
					nCho = -1
					nJong = 12
				} else if nCho == 5 && p == 17 { // ㄿ
					nCho = -1
					nJong = 13
				} else if nCho == 5 && p == 18 { // ㅀ
					nCho = -1
					nJong = 14
				} else if nCho == 7 && p == 9 { // ㅄ
					nCho = -1
					nJong = 17
				} else {							// 단자음을 연타
					res += CHO_DATA[nCho]
					nCho = indexOf(KOR_KEY[p], CHO_DATA)
				}
			}
		} else {									// 모음
			if nJong != -1 {						// (앞글자 종성), 초성+중성
				// 복자음 다시 분해
				var newCho int  // (임시용) 초성
				if nJong == 2 { // ㄱ, ㅅ
					nJong = 0
					newCho = 9
				} else if nJong == 4 { // ㄴ, ㅈ
					nJong = 3
					newCho = 12
				} else if nJong == 5 { // ㄴ, ㅎ
					nJong = 3
					newCho = 18
				} else if nJong == 8 { // ㄹ, ㄱ
					nJong = 7
					newCho = 0
				} else if nJong == 9 { // ㄹ, ㅁ
					nJong = 7
					newCho = 6
				} else if nJong == 10 { // ㄹ, ㅂ
					nJong = 7
					newCho = 7
				} else if nJong == 11 { // ㄹ, ㅅ
					nJong = 7
					newCho = 9
				} else if nJong == 12 { // ㄹ, ㅌ
					nJong = 7
					newCho = 16
				} else if nJong == 13 { // ㄹ, ㅍ
					nJong = 7
					newCho = 17
				} else if nJong == 14 { // ㄹ, ㅎ
					nJong = 7
					newCho = 18
				} else if nJong == 17 { // ㅂ, ㅅ
					nJong = 16
					newCho = 9
				} else {							// 복자음 아님
					newCho = indexOf(JONG_DATA[nJong], CHO_DATA)
					nJong = -1
				}
				if nCho != -1 { // 앞글자가 초성+중성+(종성)
					res += makeHangul(nCho, nJung, nJong)
				} else { // 복자음만 있음
					res += JONG_DATA[nJong]
				}

				nCho = newCho
				nJung = -1
				nJong = -1
			}
			if nJung == -1 { // 중성 입력 중
				nJung = indexOf(KOR_KEY[p], JUNG_DATA)
			} else if nJung == 8 && p == 19 { // ㅘ
				nJung = 9
			} else if nJung == 8 && p == 20 { // ㅙ
				nJung = 10
			} else if nJung == 8 && p == 32 { // ㅚ
				nJung = 11
			} else if nJung == 13 && p == 23 { // ㅝ
				nJung = 14
			} else if nJung == 13 && p == 24 { // ㅞ
				nJung = 15
			} else if nJung == 13 && p == 32 { // ㅟ
				nJung = 16
			} else if nJung == 18 && p == 32 { // ㅢ
				nJung = 19
			} else {            // 조합 안되는 모음 입력
				if nCho != -1 { // 초성+중성 후 중성
					res += makeHangul(nCho, nJung, nJong)
					nCho = -1
				} else	{					// 중성 후 중성
					res += JUNG_DATA[nJung]
				}
				nJung = -1
				res += string(KOR_KEY[p])
			}
		}
	}

	// 마지막 한글이 있으면 처리
	if nCho != -1 {
		if nJung != -1 { // 초성+중성+(종성)
			res += makeHangul(nCho, nJung, nJong)
		} else { // 초성만
			res += CHO_DATA[nCho]
		}
	} else {
		if nJung != -1 { // 중성만
			res += JUNG_DATA[nJung]
		} else {						// 복자음
			if nJong != -1 {
				res += JONG_DATA[nJong]
			}
		}
	}

	return res
}

func makeHangul(nCho int, nJung int, nJong int) string {
	return fmt.Sprintf("%v", string(0xac00 + nCho*21*28 + nJung*28 + nJong + 1))
}
