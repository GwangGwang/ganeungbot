package main

import (
	"fmt"
	"github.com/GwangGwang/ganeungbot/internal/pkg/db"
	"github.com/GwangGwang/ganeungbot/pkg/lol"
)

const (
	hi = iota
	hii
)

/*
assumption: chatgroup / users collections already pre-populated
1. fetch summoner infos - upsert if necessary
2. for each summoner, fetch matchlist
3. reformat matchlist and upload
  gather match ids mid-way
4. fetch matches
5. reformat matches in neat format



 */

func main() {
//	fmt.Println("connecting to mongo")
	_ = db.ConnectDB()


	db.Hi = "what"

	fmt.Println(db.Hi)
//	fmt.Println("connected to mongo!")
//
//
	key := "RGAPI-0acaf874-ffce-45fa-9a1f-696a8351a680"
	lolObj, err := lol.New(key)
	if err != nil {
		panic(err)
	}

	err = lolObj.Update()
	if err != nil {
		panic(err)
	}

}

/*
{ "_id" : ObjectId("5e0fa3aeac2eacbefaa43dbb"), "humanname" : "광승", "username" : "gwanggwang", "summonerNames" : [ "GwangGwang", "KwangKwang", "KimGwangGwang" ] }
{ "_id" : ObjectId("5e0fa4afac2eacbefaa43dbc"), "humanname" : "영하", "username" : "younghaan", "summonerNames" : [ "0ha", "1ha", "looc", "3ha", "5ha" ] }
{ "_id" : ObjectId("5e0fa4e7ac2eacbefaa43dbd"), "humanname" : "은국", "username" : "silversoup", "summonerNames" : [ "SilverSoup" ] }
{ "_id" : ObjectId("5e0fa4faac2eacbefaa43dbe"), "humanname" : "소라", "username" : "", "summonerNames" : [ "Laya Yi" ] }
{ "_id" : ObjectId("5e0fa4ffac2eacbefaa43dbf"), "humanname" : "형주", "username" : "appiejam", "summonerNames" : [ "appiejam", "LoveHeals", "LoveEndures" ] }
{ "_id" : ObjectId("5e0fa503ac2eacbefaa43dc0"), "humanname" : "찬주", "username" : "chanjook", "summonerNames" : [ "cj2da" ] }
 */
