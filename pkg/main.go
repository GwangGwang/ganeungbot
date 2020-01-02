package main

import (
	"fmt"
	"github.com/GwangGwang/ganeungbot/internal/pkg/db"
)

const (
	hi = iota
	hii
)

func main() {

	fmt.Println("connecting")
	err := db.ConnectDB()
	fmt.Println("connected!")

	if err != nil {
		panic(err)
	}

	db.InsertUserLocation("gwanggwang", "toronto")
	loc := db.GetUserLocation("gwanggwang")

	fmt.Println(loc)



}
