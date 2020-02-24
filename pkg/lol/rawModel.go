package lol

/*
 RAW DATA FROM RIOT API

 These models exist for fetching/uploading to db.
 Actual logic computations are held with better models, found in model.go
*/
type MatchlistRaw struct {
	MatchReferencesRaw []MatchReferenceRaw `json:"matches" bson:"matches"`
	TotalGames         int                 `json:"totalGames"`
	StartIndex         int                 `json:"startIndex"`
	EndIndex           int                 `json:"endIndex"`
	SummonerName       string              `json:"summonerName" bson:"summonerName"` // only used internally
}

type MatchReferenceRaw struct {
	Timestamp  int64  `json:"timestamp"`
	Lane       string `json:"lane"` // TODO: investigate what this entails
	Role       string `json:"role"` // TODO: same as above
	GameId     int64  `json:"gameId"`
	ChampionId int    `json:"champion"`
	PlatformId string `json:"platformId"`
	Season     int    `json:"season"`
	QueueId    int    `json:"queue"`
	SummonerName string // only used internally
}

type MatchRaw struct {
	SeasonId    int    `json:"seasonId"`
	QueueId     int    `json:"queueId"`
	GameId      int64  `json:"gameId"`
	GameMode    string `json:"gameMode"`
	GameType    string `json:"gameType"`
	MapId       int    `json:"mapId"`
	PlatformId  string `json:"platformId"`
	GameVersion string `json:"gameVersion"`

	GameCreation int64 `json:"gameCreation"`
	GameDuration int64 `json:"gameDuration"`

	ParticipantIdentities []ParticipantIdentity `json:"participantIdentities" bson:"participantIdentities"` // name, accountId, etc.
	TeamStats             []TeamStat            `json:"teams" bson:"teams"`
	ParticipantData       []ParticipantData     `json:"participants" bson:"participants"`
}

type ParticipantIdentity struct {
	ParticipantId int           `json:"participantId"`
	Data          PlayerDataRaw `json:"player" bson:"player"`
}

type PlayerDataRaw struct {
	SummonerName     string `json:"summonerName"`
	CurrentAccountId string `json:"currentAccountId"`
	SummonerId       string `json:"summonerId"`
	AccountId        string `json:"accountId"`
}

type TeamStat struct {
	TeamId int    `json:"teamId"` // 100 - blue / 200 - red
	Win    string `json:"win"`    // Fail / Win

	FirstBlood     bool `json:"firstBlood"`
	FirstInhibitor bool `json:"firstInhibitor"`
	FirstTower     bool `json:"firstTower"`

	InhibitorKills int `json:"inhibitorKills"`
	TowerKills     int `json:"towerKills"`

	// relevant to only Summoner's Rift
	FirstDragon     bool `json:"firstDragon"`
	FirstRiftHerald bool `json:"firstRiftHerald"`
	FirstBaron      bool `json:"firstBaron"`

	DragonKills     int `json:"dragonKills"`
	RiftHeraldKills int `json:"riftHeraldKills"`
	BaronKills      int `json:"baronKills"`

	// only draft games
	Bans []TeamBan `json:"bans" bson:"bans"`
}

type TeamBan struct {
	PickTurn   int `json:"pickTurn"`
	ChampionId int `json:"championId"`
}

type ParticipantData struct {
	ParticipantId int `json:"participantId"`
	TeamId        int `json:"teamId"`
	Spell1Id      int `json:"spell1Id"`
	Spell2Id      int `json:"spell2Id"`
	ChampionId    int `json:"championId"`

	Stats    ParticipantStats        `json:"stats" bson:"stats"`
	Timeline ParticipantTimelineData `json:"timeline" bson:"timeline"`
}

type ParticipantStats struct {
	ParticipantId int  `json:"participantId"`
	Win           bool `json:"win"`

	// KDA / Kills
	Kills   int `json:"kills"`
	Deaths  int `json:"deaths"`
	Assists int `json:"assists"`

	FirstBloodKill   bool `json:"firstBloodKill"`
	FirstBloodAssist bool `json:"firstBloodAssist"`

	KillingSprees       int `json:"killingSprees"`
	LargestKillingSpree int `json:"largestKillingSpree"`

	LargestMultiKill int `json:"largestMultiKill"`
	DoubleKills      int `json:"doubleKills"`
	TripleKills      int `json:"tripleKills"`
	QuadraKills      int `json:"quadraKills"`
	PentaKills       int `json:"pentaKills"`
	UnrealKills      int `json:"unrealKills"`

	// Damage
	TotalDamageDealtToChampions    float64 `json:"totalDamageDealtToChampions"`
	PhysicalDamageDealtToChampions float64 `json:"physicalDamageDealtToChampions"`
	MagicDamageDealtToChampions    float64 `json:"magicDamageDealtToChampions"`
	TrueDamageDealtToChampions     float64 `json:"trueDamageDealtToChampions"`

	TotalDamageDealt    float64 `json:"totalDamageDealt"`
	PhysicalDamageDealt float64 `json:"physicalDamageDealt"`
	MagicDamageDealt    float64 `json:"magicDamageDealt"`
	TrueDamageDealt     float64 `json:"trueDamageDealt"`

	LargestCriticalStrike int `json:"largestCriticalStrike"`

	// Objectives
	TeamObjective int `json:"teamObjective"`

	DamageDealtToObjectives float64 `json:"damageDealtToObjectives"`
	DamageDealtToTurrets    float64 `json:"damageDealtToTurrets"`

	TurretKills      int  `json:"turretKills"`
	FirstTowerKill   bool `json:"firstTowerKill"`
	FirstTowerAssist bool `json:"firstTowerAssist"`

	InhibitorKills       int  `json:"inhibitorKills"`
	FirstInhibitorKill   bool `json:"firstInhibitorKill"`
	FirstInhibitorAssist bool `json:"firstInhibitorAssist"`

	// Healing
	TotalHeal        float64 `json:"totalHeal"`
	TotalUnitsHealed int     `json:"totalUnitsHealed"`

	// Tanking
	TotalDamageTaken    float64 `json:"totalDamageTaken"`
	PhysicalDamageTaken float64 `json:"physicalDamageTaken"`
	MagicalDamageTaken  float64 `json:"magicalDamageTaken"`
	TrueDamageTaken     float64 `json:"trueDamageTaken"`

	DamageSelfMitigated float64 `json:"damageSelfMitigated"`

	// CS
	TotalMinionsKilled int `json:"totalMinionsKilled"`

	NeutralMinionsKilled            int `json:"neutralMinionsKilled"`
	NeutralMinionsKilledTeamJungle  int `json:"neutralMinionsKilledTeamJungle"`
	NeutralMinionsKilledEnemyJungle int `json:"neutralMinionsKilledEnemyJungle"`

	// Vision
	VisionScore             float64 `json:"visionScore"`
	VisionWardsBoughtInGame int     `json:"visionWardsBoughtInGame"`
	SightWardsBoughtInGame  int     `json:"sightWardsBoughtInGame"`

	WardsPlaced int `json:"wardsPlaced"`
	WardsKilled int `json:"wardsKilled"`

	// CC
	TimeCCingOthers            float64 `json:"timeCCingOthers"` // TODO: what is this
	TotalTimeCrowdControlDealt int     `json:"totalTimeCrowdControlDealt"`

	// etc.
	LongestTimeSpentLiving int `json:"longestTimeSpentLiving"`

	//	NodeCapture          int `json:"nodeCapture"`
	//	NodeCaptureAssist    int `json:"nodeCaptureAssist"`
	//	NodeNeutralize       int `json:"nodeNeutralize"`
	//	NodeNeutralizeAssist int `json:"nodeNeutralizeAssist"`
	//	AltarsCaptured       int `json:"altarsCaptured"`
	//	AltarsNeutralized    int `json:"altarsNeutralized"`

	TotalPlayerScore  int `json:"totalPlayerScore"`
	CombatPlayerScore int `json:"combatPlayerScore"`

	GoldEarned int `json:"goldEarned"`
	GoldSpent  int `json:"goldSpent"`

	ChampLevel int `json:"champLevel"`

	//	Item0 int `json:"item0"`
	//	Item1 int `json:"item1"`
	//	Item2 int `json:"item2"`
	//	Item3 int `json:"item3"`
	//	Item4 int `json:"item4"`
	//	Item5 int `json:"item5"`
	//	Item6 int `json:"item6"`
	//
	//	PerkPrimaryStyle int `json:"perkPrimaryStyle"`
	//	PerkSubStyle     int `json:"perkSubStyle"`
	//
	//	Perk0 int `json:"perk0"`
	//	Perk1 int `json:"perk1"`
	//	Perk2 int `json:"perk2"`
	//	Perk3 int `json:"perk3"`
	//	Perk4 int `json:"perk4"`
	//	Perk5 int `json:"perk5"`
	//
	//	Perk0Var1 int `json:"perk0Var1"`
	//	Perk0Var2 int `json:"perk0Var2"`
	//	Perk0Var3 int `json:"perk0Var3"`
	//	Perk1Var1 int `json:"perk1Var1"`
	//	Perk1Var2 int `json:"perk1Var2"`
	//	Perk1Var3 int `json:"perk1Var3"`
	//	Perk2Var1 int `json:"perk2Var1"`
	//	Perk2Var2 int `json:"perk2Var2"`
	//	Perk2Var3 int `json:"perk2Var3"`
	//	Perk3Var1 int `json:"perk3Var1"`
	//	Perk3Var2 int `json:"perk3Var2"`
	//	Perk3Var3 int `json:"perk3Var3"`
	//	Perk4Var1 int `json:"perk4Var1"`
	//	Perk4Var2 int `json:"perk4Var2"`
	//	Perk4Var3 int `json:"perk4Var3"`
	//	Perk5Var1 int `json:"perk5Var1"`
	//	Perk5Var2 int `json:"perk5Var2"`
	//	Perk5Var3 int `json:"perk5Var3"`
	//
	//	ObjectivePlayerScore int `json:"objectivePlayerScore"`
	//	TotalScoreRank       int `json:"totalScoreRank"`
	//	PlayerScore0         int `json:"playerScore0"` // what are these?
	//	PlayerScore1         int `json:"playerScore1"`
	//	PlayerScore2         int `json:"playerScore2"`
	//	PlayerScore3         int `json:"playerScore3"`
	//	PlayerScore4         int `json:"playerScore4"`
	//	PlayerScore5         int `json:"playerScore5"`
	//	PlayerScore6         int `json:"playerScore6"`
	//	PlayerScore7         int `json:"playerScore7"`
	//	PlayerScore8         int `json:"playerScore8"`
	//	PlayerScore9         int `json:"playerScore9"`
}

type ParticipantTimelineData struct {
	Lane          string `json:"lane"`
	Role          string `json:"role"`
	ParticipantId int    `json:"participantId"`

	// CS
	GoldPerMinDeltas   map[string]float64 `json:"goldPerMinDeltas" bson:"goldPerMinDeltas"`
	CreepsPerMinDeltas map[string]float64 `json:"creepsPerMinDeltas" bson:"creepsPerMinDeltas"`
	CsDiffPerMinDeltas map[string]float64 `json:"csDiffPerMinDeltas" bson:"csDiffPerMinDeltas"`

	// XP
	XpDiffPerMinDeltas map[string]float64 `json:"xpDiffPerMinDeltas" bson:"xpDiffPerMinDeltas"`
	XpPerMinDeltas     map[string]float64 `json:"xpPerMinDeltas" bson:"xpPerMinDeltas"`

	// Damage
	DamageTakenDiffPerMinDeltas map[string]float64 `json:"damageTakenDiffPerMinDeltas" bson:"damageTakenDiffPerMinDeltas"`
	DamageTakenPerMinDeltas     map[string]float64 `json:"damageTakenPerMinDeltas" bson:"damageTakenPerMinDeltas"`
}

var Champions = []ChampionInfo{
	{Id: "89", Name: "Leona", Matchers: []string{"레오나"}},
	{Id: "110", Name: "Varus", Matchers: []string{"바루스", "게이"}},
	{Id: "111", Name: "Nautilus", Matchers: []string{"노틸", "노딜러스"}},
	{Id: "112", Name: "Viktor", Matchers: []string{"빅토르"}},
	{Id: "113", Name: "Sejuani", Matchers: []string{"세주아니", "세주"}},
	{Id: "114", Name: "Fiora", Matchers: []string{"피오라"}},
	{Id: "236", Name: "Lucian", Matchers: []string{"루시안"}},
	{Id: "115", Name: "Ziggs", Matchers: []string{"직스"}},
	{Id: "117", Name: "Lulu", Matchers: []string{"룰루"}},
	{Id: "90", Name: "Malzahar", Matchers: []string{"말자하"}},
	{Id: "238", Name: "Zed", Matchers: []string{"제드"}},
	{Id: "91", Name: "Talon", Matchers: []string{"탈론"}},
	{Id: "119", Name: "Draven", Matchers: []string{"드레이븐"}},
	{Id: "92", Name: "Riven", Matchers: []string{"리븐"}},
	{Id: "516", Name: "Ornn", Matchers: []string{"오른"}},
	{Id: "96", Name: "Kog'Maw", Matchers: []string{"코그모"}},
	{Id: "10", Name: "Kayle", Matchers: []string{"케일"}},
	{Id: "98", Name: "Shen", Matchers: []string{"쉔"}},
	{Id: "99", Name: "Lux", Matchers: []string{"럭스"}},
	{Id: "11", Name: "Master Yi", Matchers: []string{"마이", "마스터이"}},
	{Id: "12", Name: "Alistar", Matchers: []string{"알리", "알리스타"}},
	{Id: "13", Name: "Ryze", Matchers: []string{"라이즈"}},
	{Id: "14", Name: "Sion", Matchers: []string{"사이온"}},
	{Id: "15", Name: "Sivir", Matchers: []string{"시비르"}},
	{Id: "16", Name: "Soraka", Matchers: []string{"소라카", "소라", "라카"}},
	{Id: "17", Name: "Teemo", Matchers: []string{"티모"}},
	{Id: "18", Name: "Tristana", Matchers: []string{"트타", "트리스타나"}},
	{Id: "19", Name: "Warwick", Matchers: []string{"워윅", "위윅", "위웍", "워웍"}},
	{Id: "240", Name: "Kled", Matchers: []string{"클레드"}},
	{Id: "120", Name: "Hecarim", Matchers: []string{"헤카림"}},
	{Id: "121", Name: "Kha'Zix", Matchers: []string{"카직스"}},
	{Id: "1", Name: "Annie", Matchers: []string{"애니", "티버"}},
	{Id: "122", Name: "Darius", Matchers: []string{"다리우스"}},
	{Id: "2", Name: "Olaf", Matchers: []string{"올라프"}},
	{Id: "245", Name: "Ekko", Matchers: []string{"에코"}},
	{Id: "3", Name: "Galio", Matchers: []string{"갈리오"}},
	{Id: "4", Name: "Twisted Fate", Matchers: []string{"트페", "트위스티드페이트"}},
	{Id: "126", Name: "Jayce", Matchers: []string{"제이스"}},
	{Id: "5", Name: "Xin Zhao", Matchers: []string{"신짜오", "짜장"}},
	{Id: "127", Name: "Lissandra", Matchers: []string{"리산드라"}},
	{Id: "6", Name: "Urgot", Matchers: []string{"우르곳"}},
	{Id: "7", Name: "LeBlanc", Matchers: []string{"르블랑"}},
	{Id: "8", Name: "Vladimir", Matchers: []string{"블라디", "모기", "블라디미르"}},
	{Id: "9", Name: "Fiddlesticks", Matchers: []string{"피들", "피들스틱", "까아악"}},
	{Id: "20", Name: "Nunu", Matchers: []string{"누누"}},
	{Id: "21", Name: "Miss Fortune", Matchers: []string{"미포", "미스포츈", "미스포춘"}},
	{Id: "22", Name: "Ashe", Matchers: []string{"애쉬"}},
	{Id: "23", Name: "Tryndamere", Matchers: []string{"트린", "트린다", "트린다미어"}},
	{Id: "24", Name: "Jax", Matchers: []string{"잭스"}},
	{Id: "25", Name: "Morgana", Matchers: []string{"몰가", "모르가나"}},
	{Id: "26", Name: "Zilean", Matchers: []string{"질리언"}},
	{Id: "27", Name: "Singed", Matchers: []string{"신지드"}},
	{Id: "28", Name: "Evelynn", Matchers: []string{"이블린", "에블린"}},
	{Id: "29", Name: "Twitch", Matchers: []string{"트위치"}},
	{Id: "131", Name: "Diana", Matchers: []string{"다이애나", "다이아나"}},
	{Id: "133", Name: "Quinn", Matchers: []string{"퀸", "까악"}},
	{Id: "254", Name: "Vi", Matchers: []string{"바이"}},
	{Id: "497", Name: "Rakan", Matchers: []string{"라칸"}},
	{Id: "134", Name: "Syndra", Matchers: []string{"신드라"}},
	{Id: "498", Name: "Xayah", Matchers: []string{"자야"}},
	{Id: "136", Name: "Aurelion Sol", Matchers: []string{"아우솔", "솔", "아우렐리언솔"}},
	{Id: "412", Name: "Thresh", Matchers: []string{"쓰레쉬"}},
	{Id: "30", Name: "Karthus", Matchers: []string{"카서스"}},
	{Id: "31", Name: "Cho'Gath", Matchers: []string{"초가스"}},
	{Id: "32", Name: "Amumu", Matchers: []string{"아무무", "무무"}},
	{Id: "33", Name: "Rammus", Matchers: []string{"람머스", "그래"}},
	{Id: "34", Name: "Anivia", Matchers: []string{"에니비아", "애니비아"}},
	{Id: "35", Name: "Shaco", Matchers: []string{"샤코"}},
	{Id: "36", Name: "Dr. Mundo", Matchers: []string{"문도"}},
	{Id: "37", Name: "Sona", Matchers: []string{"소나"}},
	{Id: "38", Name: "Kassadin", Matchers: []string{"카사딘"}},
	{Id: "39", Name: "Irelia", Matchers: []string{"이렐", "이렐리아"}},
	{Id: "141", Name: "Kayn", Matchers: []string{"케인"}},
	{Id: "142", Name: "Zoe", Matchers: []string{"조이"}},
	{Id: "143", Name: "Zyra", Matchers: []string{"자이라"}},
	{Id: "266", Name: "Aatrox", Matchers: []string{"아트럭스"}},
	{Id: "420", Name: "Illaoi", Matchers: []string{"일라오이"}},
	{Id: "145", Name: "Kai'Sa", Matchers: []string{"카이사"}},
	{Id: "267", Name: "Nami", Matchers: []string{"나미"}},
	{Id: "421", Name: "Rek'Sai", Matchers: []string{"렉사이"}},
	{Id: "268", Name: "Azir", Matchers: []string{"아지르"}},
	{Id: "427", Name: "Ivern", Matchers: []string{"아이번"}},
	{Id: "429", Name: "Kalista", Matchers: []string{"칼리스타", "칼리"}},
	{Id: "40", Name: "Janna", Matchers: []string{"잔나"}},
	{Id: "41", Name: "Gangplank", Matchers: []string{"갱플랭크", "갱플"}},
	{Id: "42", Name: "Corki", Matchers: []string{"코르키"}},
	{Id: "43", Name: "Karma", Matchers: []string{"카르마"}},
	{Id: "44", Name: "Taric", Matchers: []string{"타릭"}},
	{Id: "45", Name: "Veigar", Matchers: []string{"베이가"}},
	{Id: "48", Name: "Trundle", Matchers: []string{"트런들"}},
	{Id: "150", Name: "Gnar", Matchers: []string{"나르"}},
	{Id: "154", Name: "Zac", Matchers: []string{"자크", "젤리"}},
	{Id: "432", Name: "Bard", Matchers: []string{"바드"}},
	{Id: "157", Name: "Yasuo", Matchers: []string{"야스오"}},
	{Id: "50", Name: "Swain", Matchers: []string{"스웨인"}},
	{Id: "51", Name: "Caitlyn", Matchers: []string{"케이틀린", "케틀"}},
	{Id: "53", Name: "Blitzcrank", Matchers: []string{"블리츠크랭크", "블리츠", "블츠"}},
	{Id: "54", Name: "Malphite", Matchers: []string{"말파이트", "말파"}},
	{Id: "55", Name: "Katarina", Matchers: []string{"카타리나", "카타"}},
	{Id: "56", Name: "Nocturne", Matchers: []string{"녹턴"}},
	{Id: "57", Name: "Maokai", Matchers: []string{"마오카이", "마오"}},
	{Id: "58", Name: "Renekton", Matchers: []string{"레넥톤", "레넥턴", "레넥"}},
	{Id: "59", Name: "Jarvan IV", Matchers: []string{"자르반4세", "자르반", "자반"}},
	{Id: "161", Name: "Vel'Koz", Matchers: []string{"벨코즈"}},
	{Id: "163", Name: "Taliyah", Matchers: []string{"탈리야"}},
	{Id: "164", Name: "Camille", Matchers: []string{"카밀"}},
	{Id: "201", Name: "Braum", Matchers: []string{"브라움"}},
	{Id: "202", Name: "Jhin", Matchers: []string{"진"}},
	{Id: "203", Name: "Kindred", Matchers: []string{"킨드레드", "킨드"}},
	{Id: "60", Name: "Elise", Matchers: []string{"엘리스"}},
	{Id: "61", Name: "Orianna", Matchers: []string{"오리아나"}},
	{Id: "62", Name: "Wukong", Matchers: []string{"오공"}},
	{Id: "63", Name: "Brand", Matchers: []string{"브랜드"}},
	{Id: "64", Name: "Lee Sin", Matchers: []string{"리신"}},
	{Id: "67", Name: "Vayne", Matchers: []string{"베인"}},
	{Id: "68", Name: "Rumble", Matchers: []string{"럼블"}},
	{Id: "69", Name: "Cassiopeia", Matchers: []string{"카시오페아", "카시오피아", "카시"}},
	{Id: "72", Name: "Skarner", Matchers: []string{"스카너"}},
	{Id: "74", Name: "Heimerdinger", Matchers: []string{"하이머딩거", "하이머"}},
	{Id: "75", Name: "Nasus", Matchers: []string{"나서스", "개"}},
	{Id: "76", Name: "Nidalee", Matchers: []string{"니달리"}},
	{Id: "77", Name: "Udyr", Matchers: []string{"우디르"}},
	{Id: "78", Name: "Poppy", Matchers: []string{"뽀삐"}},
	{Id: "79", Name: "Gragas", Matchers: []string{"그라가스"}},
	{Id: "222", Name: "Jinx", Matchers: []string{"징크스", "징스"}},
	{Id: "101", Name: "Xerath", Matchers: []string{"제라스"}},
	{Id: "102", Name: "Shyvana", Matchers: []string{"쉬바나", "시바나"}},
	{Id: "223", Name: "Tahm Kench", Matchers: []string{"탐켄치"}},
	{Id: "103", Name: "Ahri", Matchers: []string{"아리"}},
	{Id: "104", Name: "Graves", Matchers: []string{"그레이브즈", "그레이브스", "그브"}},
	{Id: "105", Name: "Fizz", Matchers: []string{"피즈"}},
	{Id: "106", Name: "Volibear", Matchers: []string{"볼리베어", "볼베"}},
	{Id: "80", Name: "Pantheon", Matchers: []string{"판테온", "빵테", "빵테온"}},
	{Id: "107", Name: "Rengar", Matchers: []string{"렝가"}},
	{Id: "81", Name: "Ezreal", Matchers: []string{"이즈리얼", "이즈", "이즈내놔"}},
	{Id: "82", Name: "Mordekaiser", Matchers: []string{"모데카이저", "모데"}},
	{Id: "83", Name: "Yorick", Matchers: []string{"요릭", "고인"}},
	{Id: "84", Name: "Akali", Matchers: []string{"아칼리"}},
	{Id: "85", Name: "Kennen", Matchers: []string{"케넨"}},
	{Id: "86", Name: "Garen", Matchers: []string{"가렌"}},
	{Id: "555", Name: "Pyke", Matchers: []string{"파이크"}},
	{Id: "145", Name: "Kaisa", Matchers: []string{"카이사"}},
}
