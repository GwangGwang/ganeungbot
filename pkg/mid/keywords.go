package mid

import (
	"fmt"
	"github.com/GwangGwang/ganeungbot/pkg/typehelper"
)

type Action int

const (
	ACTION_NONE Action = iota + 1
	ACTION_ABLE
	ACTION_CARDS_AGAINST_HUMANITY
	ACTION_MOST_BYUNTAE
	ACTION_GOODBYE
	ACTION_GREETINGS
	ACTION_INJUNG
	ACTION_LOL
	ACTION_TYPEHELPER
	ACTION_SHUTUP
	ACTION_UNSHUTUP
	ACTION_VERSUS
	ACTION_WEATHER
	ACTION_WORLDCUP
	ACTION_GAMESTATS
	ACTION_GAMESTATS_CHANGEMODE
)

var Commands = []struct {
	Keywords []string
	Behavior Action
}{
	{Keywords: []string{
		fmt.Sprintf("^%s", typehelper.Trigger)},
		Behavior: ACTION_TYPEHELPER},
	{Keywords: []string{
		"able\\?$",
		"possible\\?$",
		"ㄱㄴ\\?$",
		"가능\\?$",
		"되\\?$",
		"됨\\?$",
		"돼\\?$",
		"됌\\?$"},
		Behavior: ACTION_ABLE},
	{Keywords: []string{
		"^cah$",
		"^ㅋㅇㅎ$",
		"^카어휴$"},
		Behavior: ACTION_CARDS_AGAINST_HUMANITY},
	{Keywords: []string{
		"^hi",
		"^hello",
		"^ㅎㅇ$",
		"^ㅎ2$",
		"^ㄱㅁㄴ$",
		"^하이$",
		"^하잉$",
		"^안녕$",
		"^굳모닝$"},
		Behavior: ACTION_GREETINGS},
	{Keywords: []string{
		"^bye$",
		"^adios$",
		"^ㅂㅇ$",
		"^ㅂ2",
		"^바이$",
		"^빠이$",
		"^빠잉$",
		"^사요나라$",
		"^아디오스$"},
		Behavior: ACTION_GOODBYE},
	{Keywords: []string{
		"^injung\\?$",
		"^ㅇㅈ\\?$",
		"^인정\\?$"},
		Behavior: ACTION_INJUNG},
	{Keywords: []string{
		"전적$"},
		Behavior: ACTION_LOL},
	{Keywords: []string{
		"^shutup$",
		"^ㄷㅊ$",
		"^ㄲㅈ$",
		"^닥쳐$",
		"^닥치삼$",
		"^조용히해$"},
		Behavior: ACTION_SHUTUP},
	{Keywords: []string{
		"^speak$",
		"^말해$",
		"^돌아와$"},
		Behavior: ACTION_UNSHUTUP},
	{Keywords: []string{".+ vs .+"},
		Behavior: ACTION_VERSUS},
	{Keywords: []string{"날씨$", "weather$", "forecast$"},
		Behavior: ACTION_WEATHER},
	{Keywords: []string{"^월드컵"},
		Behavior: ACTION_WORLDCUP},
	{Keywords: []string{"전적$"},
		Behavior: ACTION_GAMESTATS},
	{Keywords: []string{"^전적:"},
		Behavior: ACTION_GAMESTATS_CHANGEMODE},
}

var Answers = map[Action][]string{
	ACTION_ABLE: {
		"ㄱㄱㄱ",
		"ㄲ!!!",
		"ㄱㄱㅆ",
		"ㄴㄴ",
		"ㅈㅅ;",
		"아! 가능",
		"매우 가능",
		";;자제점",
		"ㅈ ㅏ ㅈ ㅔ",
		"매우 불가능",
	},
	ACTION_GREETINGS: {
		"Hello",
		"Hi",
		"ㅎㅇ",
		"ㅎ2",
		"안녕?",
		"안녕!",
		"하잉",
		"하이",
		"헬로",
		"곤니찌와",
	},
	ACTION_GOODBYE: {
		"Bye~",
		"ㅂ2",
		"ㅂㅇ",
		"굳바이",
		"빠이",
		"빠잉",
		"사요나라",
		"아디오스!",
		"잘가요..",
	},
	ACTION_INJUNG: {
		"ㅇㅈ",
		"ㄴㄴ",
		"인정",
		"매우 인정",
	},
	ACTION_SHUTUP: {
		"당신이 원하신다면.. 흑 ㅠㅠ",
		"잘 있어요",
		"사요나라!",
		"아디오스!",
		"네 ㅠㅠㅠㅠ",
		"알았다 꺼져주지",
	},
	ACTION_UNSHUTUP: {
		"내가 돌아왔다!",
		"살았다 꺄르륵!",
		"보고싶었어요",
		"고마워요",
		"데헷 다시 떠들어야징",
		"I have returned",
	},
}
