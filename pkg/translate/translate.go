package translate

import (
	"cloud.google.com/go/translate"
	"context"
	"html"
	"fmt"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
	"log"
	"strings"
)

const Trigger = "번역:"
const usage string = `Usage:
	번역: <번역대상> - 기능 사용
	번역언어 목록
	번역언어 - 현재 설정된 언어
	번역언어 <언어코드> - 언어 변경`

// Translate is the translate object
type Translate struct {
	ApiKey   string
	Operational bool
	TargetLanguageMap map[int64]language.Tag // chatID to language
	LanguageList string
}

// New initializes and returns a new weather pkg Weather
func New(apiKey string) (Translate, error) {
	log.Println("Initializing translate pkg")

	t := Translate{}

	if len(apiKey) == 0 {
		log.Printf("WARN: weather API key not found")
		t.Operational = false
	} else {
		t.Operational = true
	}

	t.ApiKey = apiKey

	langList := make([]string, 0)
	for lang, _ := range languages {
		langList = append(langList, lang)
	}
	t.LanguageList = strings.Join(langList, ", ")

	t.TargetLanguageMap = make(map[int64]language.Tag)

	return t, nil
}

// GetResponse is the main outward facing function to generate weather response
func (t *Translate) GetResponse(chatID int64, txt string) (string, error) {
	// anonymous helper function for error returning
	handleError := func(err error) (string, error) {
		return fmt.Sprintf("%s\n%s", err.Error(), usage), err
	}

	if _, exists := t.TargetLanguageMap[chatID]; !exists {
		log.Printf("setting default language to english for newfound chatgroup '%d'", chatID)
		t.TargetLanguageMap[chatID] = language.Make("en")
	}

	parseResult, err := parse(txt)
	if err != nil {
		return "", err
	}

	switch parseResult.queryType {
	case listLanguages:
		return fmt.Sprintf("Supported langauge codes: %s", t.LanguageList), nil
	case setLanguage:
		setLang := parseResult.text
		err := t.setLanguage(chatID, setLang)
		if err != nil {
			return handleError(err)
		} else {
			langName := t.getTargetLanguage(chatID)
			return fmt.Sprintf("Changed target language to %s", langName), nil
		}
	case getLanguage:
		return fmt.Sprintf("Target language: %s", t.getTargetLanguage(chatID)), nil
	case translateText:
		resp, err := t.translate(t.TargetLanguageMap[chatID], parseResult.text)
		if err != nil {
			return handleError(fmt.Errorf("error during translation"))
		}
		return fmt.Sprintf("(%s --> %s)\n%s", getLanguageNameFromTag(resp.Source), getLanguageNameFromTag(t.TargetLanguageMap[chatID]), html.UnescapeString(resp.Text)), nil
	}

	return "", err
}

func (t *Translate) setLanguage(chatID int64, target string) error {
	if _, ok := languages[target]; ok {
		parsed, err := language.Parse(target)
		if err != nil {
			return err
		}
		t.TargetLanguageMap[chatID] = parsed
	} else {
		return fmt.Errorf("unknown language '%s'", target)
	}

	return nil
}

func (t *Translate) getTargetLanguage(chatID int64) string {
	return getLanguageNameFromTag(t.TargetLanguageMap[chatID])
}

func (t *Translate) translate(lang language.Tag, txt string) (translate.Translation, error) {
	ctx := context.Background()

	client, err := translate.NewClient(ctx, option.WithAPIKey(t.ApiKey))
	if err != nil {
		log.Print(err)
		return translate.Translation{}, err
	}
	defer client.Close()

	resp, err := client.Translate(ctx, []string{txt}, lang, nil)
	if err != nil {
		log.Print(err)
		return translate.Translation{}, err
	}

	return resp[0], nil
}

// HELPERS
func getLanguageNameFromTag(tag language.Tag) string {
	codeStr := fmt.Sprintf("%s", tag)
	if lang, ok := languages[codeStr]; ok {
		return lang
	}
	return ""
}

