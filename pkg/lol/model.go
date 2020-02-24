package lol

type queryKey struct {
	target target      // mandatory; whose stats? (ign, person, 'all')
	gameMode gameMode  // mandatory; game mode (aram, normal, etc.)
	isBestStats bool   // optional; best stats?
	championId string  // optional; which champion?
}

// aggregated data
type aggData struct {
	data map[string]float64
	resp string
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
//	{gameMode: gameModeRift, queueIds: []int{2, 4, 6, 400, 410, 420, 430, 440}, matchers: []string{"rift", "협곡"}},
	{gameMode: gameModeBot, queueIds: []int{7, 31, 32, 33, 830, 840, 850}, matchers: []string{"ai", "bot", "bots", "봇", "봇겜", "봇전"}},
}

func getQueueFromId(id int) gameMode {
	for _, queue := range queues {
		for _, queueId := range queue.queueIds {
			if queueId == id {
				return queue.gameMode
			}
		}
	}

	return gameModeNone
}



type ChampionDataRaw struct {
	Version string `json:"version"`
	Data map[string]ChampionInfo `json:"data"`
}

// Static champion info
// Base data comes from en_US api call
// Other language names and nicknames and such are populated under matches
// TODO: figure out a decent way to populate nicknames
type ChampionInfo struct {
	Id string `json:"id"`
	Key string `json:"key"`
	Name string `json:"name"`
	Matchers []string `json:"matchers" bson:"matchers"`
	Tags []string `json:"tags" bson:"tags"` // TODO: make into tag object
}

type ChatGroup struct {
	ChatID int64 `json:"chatID"`
	Users []string `json:"users" bson:"users"`
	// TODO: avg, best
}

type User struct {
	UserName string `json:"username" bson:"username"`
	HumanName string `json:"humanname" bson:"humanname"`
	SummonerNames []string `json:"summonerNames" bson:"summonerNames"`
}

func getUserFromSummoner(users []User, inputSummName string) string {
	for _, user := range users {
		for _, summName := range user.SummonerNames {
			if summName == inputSummName {
				return user.UserName
			}
		}
	}

	return ""
}

type Summoner struct {
	Name string `json:"name"`
	Level int `json:"summonerLevel"`
	RevisionDate int64 `json:"revisionDate" bson:"revisionDate"`
	Id string `json:"id"`
	AccountId string `json:"accountId"`
}

