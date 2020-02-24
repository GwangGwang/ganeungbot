package lol

import (
	"fmt"
	"github.com/GwangGwang/ganeungbot/internal/pkg/db"
	"gopkg.in/mgo.v2/bson"
	"log"
)

const (
	lolDatabase = "lol"
	scraperDatabase = "lolscraper"
)

/* STATIC DATA */
func GetUsers() []User {
	sessionCopy := db.Session.Copy()
	defer sessionCopy.Close()

	var userInfos []User
	err := sessionCopy.DB(lolDatabase).C("users").Find(bson.M{}).All(&userInfos)
	if err != nil {
		log.Printf(err.Error())
	}

	log.Printf("Retrieved all user info from db")
	return userInfos
}

func UpsertSummoners(summoners []Summoner) error {
	log.Printf("upserting summoners info")
	sessionCopy := db.Session.Copy()
	defer sessionCopy.Close()

	for _, summonerInfo := range summoners {
		query := bson.M{
			"id": summonerInfo.Id,
		}

		_, err := sessionCopy.DB(lolDatabase).C("summoners").Upsert(query, summonerInfo)
			if err != nil{
			return err
		}
	}

	return nil
}

func UpsertMatchlistRaw(matchlist MatchlistRaw) error {
	log.Printf("upserting matchlist data for summoner '%s'\n", matchlist.SummonerName)
	sessionCopy := db.Session.Copy()
	defer sessionCopy.Close()

	query := bson.M{
		"summonerName": matchlist.SummonerName,
	}

	_, err := sessionCopy.DB(lolDatabase).C("matchlistRaw").Upsert(query, matchlist)
	if err != nil {
		return fmt.Errorf("error while upserting matchlist data: %s", err.Error())
	}

	return nil
}

func UpsertMatchRaw(match MatchRaw) error {
	sessionCopy := db.Session.Copy()
	defer sessionCopy.Close()

	query := bson.M{
		"gameId": match.GameId,
	}

	_, err := sessionCopy.DB(lolDatabase).C("matchRaw").Upsert(query, match)
	if err != nil {
		return fmt.Errorf("error while upserting match id '%d': %s", match.GameId, err.Error())
	}

	return nil
}

func GetMatchRaw(gameId int64) (MatchRaw, error) {
	sessionCopy := db.Session.Copy()
	defer sessionCopy.Close()

	query := bson.M{
		"gameId": gameId,
	}

	var match MatchRaw
	if err := sessionCopy.DB(lolDatabase).C("matchRaw").Find(query).One(&match); err != nil {
		return match, err
	}

	return match, nil
}

/* Static Data */

func UpsertStaticChampionData(chdata ChampionDataRaw) {
	sessionCopy := db.Session.Copy()
	defer sessionCopy.Close()

	for _, chInfo := range chdata.Data {
		query := bson.M{
			"name": chInfo.Id,
		}

		log.Printf("inserting static data for champion %s\n", chInfo.Id)
		_, err := sessionCopy.DB(lolDatabase).C("champions").Upsert(query, chInfo)
		if err != nil {
			log.Printf(err.Error())
		}
	}
}

