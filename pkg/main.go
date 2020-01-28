package main

import (
	"fmt"
	"github.com/GwangGwang/ganeungbot/internal/pkg/db"
	"github.com/GwangGwang/ganeungbot/pkg/lolScraper"
	"log"
)

const (
	hi = iota
	hii
)

func main() {


//	fmt.Println("connecting to mongo")
	_ = db.ConnectDB()


	db.Hi = "what"

	fmt.Println(db.Hi)
//	fmt.Println("connected to mongo!")
//
//
	lol, err := lolScraper.New("RGAPI-2e16a30c-3274-4f4a-a16d-f389352b15d6")
	if err != nil {
		panic(err)
	}

	log.Printf("%+v", lol.UserInfos)


	if err != nil {
		panic(err)
	}
//
//	//lol.UpdateStaticChampionData()
//
//	wha := fmt.Sprintf("%sblahblah%d", "bleh")
//	fmt.Println(wha)
//
//	fmt.Printf(wha, 12)

}
