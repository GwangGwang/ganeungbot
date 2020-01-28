package lol

import (
	"github.com/GwangGwang/ganeungbot/internal/pkg/db"
	"gopkg.in/mgo.v2/bson"
	"log"
)

const (
	Database = "lol"
)

/* STATIC DATA */
func GetUsers() []UserInfo {
	sessionCopy := db.Session.Copy()
	defer sessionCopy.Close()


	var userInfos []UserInfo
	err := sessionCopy.DB(Database).C("users").Find(bson.M{}).All(&userInfos)
	if err != nil {
		log.Printf(err.Error())
	}

	log.Printf("Retrieved all user info from db")
	return userInfos
}


