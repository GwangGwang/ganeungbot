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
	gameModeNone gameMode = ""
	gameModeAram  = "aram"
	gameModeNormal = "normal"
	gameModeRanked = "ranked"
	gameModeRift = "rift"
	gameModeBot = "bot"
)
type queue struct {
	gameMode gameMode
	queueIds []int
	matchers []string
}

var queues = []queue{
	{gameMode: gameModeAram, queueIds: []int{65, 100, 450}, matchers: []string{"aram", "아람", "칼바람"}},
	{gameMode: gameModeNormal, queueIds: []int{2, 400, 430}, matchers: []string{"normal", "노멀", "노말", "일반"}},
	{gameMode: gameModeRanked, queueIds: []int{4, 6, 410, 420, 440}, matchers: []string{"ranked", "rank", "랭", "랭크", "솔랭"}},
	{gameMode: gameModeRift, queueIds: []int{2, 4, 6, 400, 410, 420, 430, 440}, matchers: []string{"rift", "협곡"}},
	{gameMode: gameModeBot, queueIds: []int{7, 31, 32, 33, 830, 840, 850}, matchers: []string{"ai", "bot", "bots", "봇", "봇겜", "봇전"}},
}

type championInfo struct {
	id int
	name string
	matchers []string
}

var champions = []championInfo {
	{id: 89, name: "Leona", matchers: []string{"레오나"}},
}

var championNil = championInfo{
	id: -1, name: "",
}

type UserInfo struct {
	HumanName string `json:"humanname"`
	UserName string `json:"username"`
	SummonerNames []string `json:"summonerNames" bson:"summonerNames"`
}

// Scraper related
type ChampionData struct {
	Version string `json:"version"`
	Data map[string]ChampionInfo `json:"data"`
}

type ChampionInfo struct {
	Id string `json:"id"`
	Key string `json:"key"`
	Name string `json:"name"`
	Tags []string `json:"tags"` // TODO: make into tag object
}

type SummonerInfo struct {
	Name string `json:"name"`
	Level int `json:"summonerLevel"`
	RevisionDate int64 `json:"revisionDate"`
	Id string `json:"id"`
	AccountId string `json:"accountId"`
}



