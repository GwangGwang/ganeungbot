package lol

import (
	"log"
	"github.com/GwangGwang/ganeungbot/internal/pkg/db"
	"gopkg.in/mgo.v2/bson"
)

const (
	lolDatabase = "lol"
)

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

