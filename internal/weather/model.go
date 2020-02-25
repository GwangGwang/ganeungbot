package weather

type tokenType string

const (
	none       = ""
	now        = "currently"
	day        = "daily"    // 오늘, 내일, 모레
	hour       = "hourly"   // 1시
	ampm       = "ampm"     // 오전, 오후
	hourAmpm   = "hourAmpm" // 1am, 1pm
)

type tokenGroup struct {
	Keywords []string
	Value    int64
}

type tokenInfo struct {
	NextTokens  []tokenType
	TokenGroups []tokenGroup
}

/*
 Time Keywords Lexical Syntactic Parsing
  start -> 현재, 오늘, 오전, 1시, 1pm, end
  현재 -> end
  오늘 -> 오전, 1시, 1pm, end
  오전 -> 1시, end
  1시 -> end
  1pm -> end

 Possible combinations
  currently - 현재
  daily     - 오늘
  hourly    - 오전, 1pm, 1시, 오전 1시, 오늘 1pm, 오늘 1시, 오늘 오전 1시, 오늘 오전
*/

var allTokenTypes = []tokenType{now, day, ampm, hour, hourAmpm}
var timeTokenMap = map[tokenType]tokenInfo{
	now: {
		TokenGroups: []tokenGroup{
			{Keywords: []string{"지금", "현재", "currently", "current", "now"}}},
		NextTokens: []tokenType{},
	},
	day: {
		TokenGroups: []tokenGroup{
			{Keywords: []string{"오늘", "today"}, Value: 0},
			{Keywords: []string{"내일", "tomorrow"}, Value: 1},
			{Keywords: []string{"모레"}, Value: 2},
			{Keywords: []string{"글피"}, Value: 3},
		},
		NextTokens: []tokenType{hour, ampm, hourAmpm},
	},
	ampm: {
		TokenGroups: []tokenGroup{
			{Keywords: []string{"오전"}, Value: 0},
			{Keywords: []string{"오후"}, Value: 1},
		},
		NextTokens: []tokenType{hour},
	},
	hour: {
		TokenGroups: []tokenGroup{
			{Keywords: []string{"\\d{1,2}시"}},
		},
		NextTokens: []tokenType{},
	},
	hourAmpm: {
		TokenGroups: []tokenGroup{
			{Keywords: []string{"\\d{1,2}am"}, Value: 0},
			{Keywords: []string{"\\d{1,2}pm"}, Value: 1},
		},
		NextTokens: []tokenType{},
	},
}

type locationParseGroup struct {
	Keywords []string
	Value    string
}

var locationParseGroups = []locationParseGroup{
	{Keywords: []string{"토론토", "톤토", "toronto", "Toronto ON"},
		Value: "Toronto"},
	{Keywords: []string{"밴쿠버", "벤쿠버", "vancouver"},
		Value: "Vancouver"},
}

type parseResult struct {
	TimeInfo timeinfo
	Location string
}

type timeinfo struct {
	Category tokenType
	Offset   int64
}

type token struct {
	word string
	tokenType tokenType
	value int64
}