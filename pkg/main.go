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

func main() {

	lol, err := lol.New("RGAPI-ffeadd4d-23f4-40a6-a915-c40f34898af1")
	if err != nil {
		panic(err)
	}

	fmt.Println("connecting to mongo")
	err = db.ConnectDB()
	fmt.Println("connected to mongo!")

	if err != nil {
		panic(err)
	}

	//lol.UpdateStaticChampionData()
	lol.UpdateSummonerInfo()

}
