package lolScraper

import (
	"github.com/GwangGwang/ganeungbot/internal/pkg/db"
	"github.com/GwangGwang/ganeungbot/pkg/lol"
	"gopkg.in/mgo.v2/bson"
	"log"
)

const (
	lolDatabase = lol.Database
	LOLScraperDatabase = "lolScraper"
)

func UpsertSummonerInfo(summonerInfo lol.SummonerInfo) error {
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

func UpsertMatchlist(summonerInfo lol.SummonerInfo) error {
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

func UpsertStaticChampionInfo(chinfo lol.ChampionInfo) {
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

