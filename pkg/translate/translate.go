package translate

import (
	"cloud.google.com/go/translate"
	"context"
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
	TargetLanguage language.Tag
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

	t.TargetLanguage = language.Make("English") // default is English

	return t, nil
}

// GetResponse is the main outward facing function to generate weather response
func (t *Translate) GetResponse(txt string) (string, error) {
	// anonymous helper function for error returning
	handleError := func(err error) (string, error) {
		return fmt.Sprintf("%s\n%s", err.Error(), usage), err
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
		err := t.setLanguage(setLang)
		if err != nil {
			return handleError(err)
		} else {
			langName := t.getLanguage(setLang)
			return fmt.Sprintf("Changed target language to %s", langName), nil
		}
	case getLanguage:
		return fmt.Sprintf("Target language: %s", t.getLanguage(parseResult.text)), nil
	case translateText:
		resp, err := t.translate(parseResult.text)
		if err != nil {
			return handleError(fmt.Errorf("error during translation"))
		}
		return fmt.Sprintf("(%s --> %s)\n%s", t.getLanguageFromTag(resp.Source), t.getLanguageFromTag(t.TargetLanguage), resp.Text), nil
	}

	return "", err
}

func (t *Translate) setLanguage(target string) error {
	lang, err := language.Parse(target)
	if err != nil {
		return err
	}
	t.TargetLanguage = lang
	return nil
}

func (t *Translate) getLanguage(code string) string {
	ret := code
	if lang, ok := languages[code]; ok {
		ret = lang
	}
	return ret
}

func (t *Translate) getLanguageFromTag(tag language.Tag) string {
	ret := fmt.Sprintf("%s", tag)
	return t.getLanguage(ret)
}

func (t *Translate) translate(txt string) (translate.Translation, error) {
	ctx := context.Background()

	client, err := translate.NewClient(ctx, option.WithAPIKey(t.ApiKey))
	if err != nil {
		log.Print(err)
		return translate.Translation{}, err
	}
	defer client.Close()

	resp, err := client.Translate(ctx, []string{txt}, t.TargetLanguage, nil)
	if err != nil {
		log.Print(err)
		return translate.Translation{}, err
	}

	return resp[0], nil
}
