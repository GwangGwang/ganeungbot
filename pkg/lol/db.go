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
func GetUsers() []UserInfo {
	sessionCopy := db.Session.Copy()
	defer sessionCopy.Close()


	var userInfos []UserInfo
	err := sessionCopy.DB(lolDatabase).C("users").Find(bson.M{}).All(&userInfos)
	if err != nil {
		log.Printf(err.Error())
	}

	log.Printf("Retrieved all user info from db")
	return userInfos
}


func UpsertSummonerInfo(summonerInfo SummonerInfo) error {
	sessionCopy := db.Session.Copy()
	defer sessionCopy.Close()

	query := bson.M{
		"id": summonerInfo.Id,
	}

	_, err := sessionCopy.DB(lolDatabase).C("summonerInfo").Upsert(query, summonerInfo)
	if err != nil {
		return err
	}

	return nil
}

func UpsertMatchlist(summonerInfo SummonerInfo) error {
	sessionCopy := db.Session.Copy()
	defer sessionCopy.Close()

	query := bson.M{
		"id": summonerInfo.Id,
	}

	_, err := sessionCopy.DB(lolDatabase).C("summonerInfo").Upsert(query, summonerInfo)
	if err != nil {
		return err
	}

	return nil
}

/* Static Data */

func UpsertStaticChampionInfo(chinfo ChampionInfo) {
	sessionCopy := db.Session.Copy()
	defer sessionCopy.Close()

	query := bson.M{
		"name": chinfo.Id,
	}

	log.Printf("inserting static data for champion %s\n", chinfo.Id)
	_, err := sessionCopy.DB(lolDatabase).C("championInfo").Upsert(query, chinfo)
	if err != nil {
		log.Printf(err.Error())
	}
}

