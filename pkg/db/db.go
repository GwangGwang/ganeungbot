package db

import (
	"github.com/globalsign/mgo"
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
	Hi string
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
