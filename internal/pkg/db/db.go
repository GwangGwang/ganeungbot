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
	Session *mgo.Session
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

	Session, err = mgo.DialWithInfo(info)
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


func GetUserLocation(user string) string {
	sessionCopy := Session.Copy()
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