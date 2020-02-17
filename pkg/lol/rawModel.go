package lol

/*
 * RAW DATA FROM RIOT API
 */
type MatchlistRaw struct {
	MatchReferencesRaw []MatchReferenceRaw `json:"matches" bson:"matches"`
	TotalGames int `json:"totalGames"`
	StartIndex int `json:"startIndex"`
	EndIndex int `json:"endIndex"`
}

type MatchReferenceRaw struct {
	Timestamp int64 `json:"timestamp"`
	Lane string `json:"lane"` // TODO: investigate what this entails
	Role string `json:"role"` // TODO: same as above
	GameId int64 `json:"gameId"`
	ChampionId int `json:"champion"`
	PlatformId string `json:"platformId"`
	Season int `json:"season"`
	QueueId int `json:"queue"`
}

type MatchRaw struct {
	SeasonId int `json:"seasonId"`
	QueueId int `json:"queueId"`
	GameId int64 `json:"gameId"`
	GameMode string `json:"gameMode"`
	GameType string `json:"gameType"`
	MapId int `json:"mapId"`
	PlatformId string `json:"platformId"`
	GameVersion string `json:"gameVersion"`

	GameCreation int64 `json:"gameCreation"`
	GameDuration int64 `json:"gameDuration"`

	ParticipantIdentities []ParticipantIdentity `json:"participantIdentities" bson:"participantIdentities"` // name, accountId, etc.
	TeamStats []TeamStat `json:"teams" bson:"teams"`
	ParticipantData []ParticipantData `json:"participants" bson:"participants"`
}

type ParticipantIdentity struct {
	ParticipantId int `json:"participantId"`
	Data PlayerDataRaw `json:"player" bson:"player"`
}

type PlayerDataRaw struct {
	SummonerName string `json:"summonerName"`
	CurrentAccountId string `json:"currentAccountId"`
	SummonerId string `json:"summonerId"`
	AccountId string `json:"accountId"`
}

type TeamStat struct {
	TeamId int `json:"teamId"` // 100 - blue / 200 - red
	Win string `json:"win"` // Fail / Win

	FirstBlood bool `json:"firstBlood"`
	FirstInhibitor bool `json:"firstInhibitor"`
	FirstTower bool `json:"firstTower"`

	InhibitorKills int `json:"inhibitorKills"`
	TowerKills int `json:"towerKills"`

	// relevant to only Summoner's Rift
	FirstDragon bool `json:"firstDragon"`
	FirstRiftHerald bool `json:"firstRiftHerald"`
	FirstBaron bool `json:"firstBaron"`

	DragonKills int `json:"dragonKills"`
	RiftHeraldKills int `json:"riftHeraldKills"`
	BaronKills int `json:"baronKills"`

	// only draft games
	Bans []TeamBan `json:"bans" bson:"bans"`
}

type TeamBan struct {
	PickTurn int `json:"pickTurn"`
	ChampionId int `json:"championId"`
}

type ParticipantData struct {
	ParticipantId int `json:"participantId"`
	TeamId int `json:"teamId"`
	Spell1Id int `json:"spell1Id"`
	Spell2Id int `json:"spell2Id"`
	ChampionId int `json:"championId"`

	Stats ParticipantStats `json:"stats" bson:"stats"`
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
	Lane string `json:"lane"`
	Role string `json:"role"`
	ParticipantId int `json:"participantId"`

	// CS
	GoldPerMinDeltas map[string]float64 `json:"goldPerMinDeltas" bson:"goldPerMinDeltas"`
	CreepsPerMinDeltas map[string]float64 `json:"creepsPerMinDeltas" bson:"creepsPerMinDeltas"`
	CsDiffPerMinDeltas map[string]float64 `json:"csDiffPerMinDeltas" bson:"csDiffPerMinDeltas"`

	// XP
	XpDiffPerMinDeltas map[string]float64 `json:"xpDiffPerMinDeltas" bson:"xpDiffPerMinDeltas"`
	XpPerMinDeltas map[string]float64 `json:"xpPerMinDeltas" bson:"xpPerMinDeltas"`

	// Damage
	DamageTakenDiffPerMinDeltas map[string]float64 `json:"damageTakenDiffPerMinDeltas" bson:"damageTakenDiffPerMinDeltas"`
	DamageTakenPerMinDeltas map[string]float64 `json:"damageTakenPerMinDeltas" bson:"damageTakenPerMinDeltas"`
}
