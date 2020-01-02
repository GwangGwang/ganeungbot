package db

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	hosts = "localhost:27017"
	database = "ganeungbot"
	username = ""
	password = ""
)

var (
	session *mgo.Session
)

func ConnectDB() error {
	var err error

	info := &mgo.DialInfo{
		Addrs: []string{hosts},
		Timeout: 60 * time.Second,
		Database: database,
		Username: username,
		Password: password,
	}

	session, err = mgo.DialWithInfo(info)
	if err != nil {
		return err
	}

	return nil
}

type UserLocation struct {
	ChatID string `json:"chatID"`
	Username string `json:"username"`
	Location string `json:"location"`
}

func InsertUserLocation(user string, location string) {
	sessionCopy := session.Copy()
	defer sessionCopy.Close()

	query := bson.M{"chatID": "haha", "username": user, "location": location}

	fmt.Println("inserting stuff")
	err := sessionCopy.DB(database).C("userLocation").Insert(query)
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func GetUserLocation(user string) string {
	sessionCopy := session.Copy()
	defer sessionCopy.Close()

	query := bson.M{"chatID": "haha", "username": user}

	fmt.Println("getting stuff")
	ul := UserLocation{}
	err := sessionCopy.DB(database).C("userLocation").Find(query).One(&ul)
	if err != nil {
		fmt.Printf(err.Error())
	}

	fmt.Printf("%+v", ul)
	return ul.Location
}