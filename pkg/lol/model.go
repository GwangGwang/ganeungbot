package lol

type parseResult struct {
	target target      // mandatory; whose stats? (ign, person, 'all')
	gameMode gameMode  // mandatory; game mode (aram, normal, etc.)
	isBestStats bool   // optional; best stats?
	championId int     // optional; which champion?
}

type target struct {
	category targetCategory
	name string
}

type targetCategory int
const (
	categoryNone targetCategory = iota
	categoryIgn
	categoryUser
	categoryAll
)

type gameMode string
const (
	none gameMode = ""
	aram  = "aram"
	normal = "normal"
	ranked = "ranked"
	rift = "rift"
	bot = "bot"
)
type queue struct {
	gameMode gameMode
	queueIds []int
	matchers []string
}

var queues = []queue{
	{gameMode: aram, queueIds: []int{65, 100, 450}, matchers: []string{"aram", "아람", "칼바람"}},
	{gameMode: normal, queueIds: []int{2, 400, 430}, matchers: []string{"normal", "노멀", "노말", "일반"}},
	{gameMode: ranked, queueIds: []int{4, 6, 410, 420, 440}, matchers: []string{"ranked", "rank", "랭", "랭크", "솔랭"}},
	{gameMode: rift, queueIds: []int{2, 4, 6, 400, 410, 420, 430, 440}, matchers: []string{"rift", "협곡"}},
	{gameMode: bot, queueIds: []int{7, 31, 32, 33, 830, 840, 850}, matchers: []string{"ai", "bot", "bots", "봇", "봇겜", "봇전"}},
}

var usermap = map[string]userinfo{
	"광승": {igns: []string{"GwangGwang", "KwangKwang"}},
	"영하": {igns: []string{"0ha", "1ha", "looc", "3ha", "5ha"}},
	"은국": {igns: []string{"SilverSoup"}},
	"찬주": {igns: []string{"cj2da"}},
	"형주": {igns: []string{"appiejam", "LoveHeals", "LoveEndures"}},
	"소라": {igns: []string{"Laya Yi"}},
}

type userinfo struct {
	igns []string
}



