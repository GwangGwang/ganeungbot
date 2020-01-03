package lol

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)





/* STATIC DATA */
func (l *LOL) UpdateStaticChampionData() error {
	url := "http://ddragon.leagueoflegends.com/cdn/9.24.2/data/en_US/champion.json"
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error while retrieving champion data: %s", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var chdata ChampionData
	err = json.Unmarshal(body, &chdata)
	if err != nil {
		return fmt.Errorf("error while unmarshalling champion data json body: %s", err)
	}

	for _, chInfo := range chdata.Data {
		UpsertStaticChampionInfo(chInfo)
	}

	return nil
}
