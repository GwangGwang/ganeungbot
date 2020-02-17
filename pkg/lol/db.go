package lol

import (
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

func UpsertMatchlist(summonerInfo Summoner) error {
	sessionCopy := db.Session.Copy()
	defer sessionCopy.Close()

	query := bson.M{
		"id": summonerInfo.Id,
	}

	_, err := sessionCopy.DB(lolDatabase).C("summoners").Upsert(query, summonerInfo)
	if err != nil {
		return err
	}

	return nil
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

