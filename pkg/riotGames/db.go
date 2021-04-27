package riotGames

import (
	"fmt"
	"github.com/GwangGwang/ganeungbot/pkg/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

const (
	database = "lol"

	championsCollection = "champions"
	usersCollection = "users"
	summonersCollection = "summoners"
	matchlistsCollection = "matchlists"
	matchCollection = "matches"
)

func ensureIndexes() error {
	if db.Session == nil {
		conn := db.ConnectDB()
		if conn != nil {
			log.Fatal("Error on connecting to MongoDB database")
			return fmt.Errorf("could not connect to MongoDB database")
		}
	}

	sessionCopy := db.Session.Copy()
	defer sessionCopy.Close()

	collectionToIndex := map[string][]string{
		championsCollection: {"id"},
		usersCollection: {"username"},
		summonersCollection: {"name"},
		matchlistsCollection: {"summonerName"},
		matchCollection: {"gameId"},
	}

	for collection, indices := range collectionToIndex {
		index := mgo.Index{
			Key: indices,
			Unique:     true,
			DropDups:   false,
			Background: true,
			Sparse:     false,
		}

		if err := sessionCopy.DB(database).C(collection).EnsureIndex(index); err != nil {
			errstr := fmt.Sprintf("error on creating %v index in the %s collection: %s", indices, collection, err)
			log.Fatal(errstr)
			return fmt.Errorf(errstr)
		}

	}

	log.Print("successfully created all db indices")
	return nil
}


/* STATIC DATA */
func GetUsers() []User {
	sessionCopy := db.Session.Copy()
	defer sessionCopy.Close()

	var userInfos []User
	err := sessionCopy.DB(database).C(usersCollection).Find(bson.M{}).All(&userInfos)
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

		_, err := sessionCopy.DB(database).C(summonersCollection).Upsert(query, summonerInfo)
		if err != nil{
			return err
		}
	}

	return nil
}

func UpsertMatchlist(matchlist Matchlist) error {
	log.Printf("upserting matchlist data for summoner '%s'\n", matchlist.SummonerName)
	sessionCopy := db.Session.Copy()
	defer sessionCopy.Close()

	query := bson.M{
		"summonerName": matchlist.SummonerName,
	}

	_, err := sessionCopy.DB(database).C(matchlistsCollection).Upsert(query, matchlist)
	if err != nil {
		return fmt.Errorf("error while upserting matchlist data: %s", err.Error())
	}

	return nil
}

func UpsertMatch(match Match) error {
	sessionCopy := db.Session.Copy()
	defer sessionCopy.Close()

	query := bson.M{
		"gameId": match.GameId,
	}

	_, err := sessionCopy.DB(database).C(matchCollection).Upsert(query, match)
	if err != nil {
		return fmt.Errorf("error while upserting match id '%d': %s", match.GameId, err.Error())
	}

	return nil
}

func GetMatchlist(summoner Summoner) (Matchlist, error) {
	sessionCopy := db.Session.Copy()
	defer sessionCopy.Close()

	query := bson.M{
		"summonerName": summoner.Name,
	}

	var matchlist Matchlist
	if err := sessionCopy.DB(database).C(matchlistsCollection).Find(query).One(&matchlist); err != nil {
		return matchlist, err
	}

	return matchlist, nil
}

func GetMatchRaw(gameId int64) (Match, error) {
	sessionCopy := db.Session.Copy()
	defer sessionCopy.Close()

	query := bson.M{
		"gameId": gameId,
	}

	var match Match
	if err := sessionCopy.DB(database).C(matchCollection).Find(query).One(&match); err != nil {
		return match, err
	}

	return match, nil
}

/* Static AggStats */

func UpsertStaticChampionData(chdata ChampionData) {
	sessionCopy := db.Session.Copy()
	defer sessionCopy.Close()

	for _, chInfo := range chdata.Data {
		query := bson.M{
			"name": chInfo.Id,
		}

		log.Printf("inserting static data for champion %s\n", chInfo.Id)
		_, err := sessionCopy.DB(database).C(championsCollection).Upsert(query, chInfo)
		if err != nil {
			log.Printf(err.Error())
		}
	}
}

