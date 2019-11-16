package main

import (
	"fmt"

	"github.com/GwangGwang/ganeungbot/pkg/weather"
)

func main() {
	fmt.Println("test")

	w, _ := weather.New()

	ll, _ := w.GetGPS("vancouver")

	fmt.Printf("%+v", ll)

	//	fmt.Printf("results: %+v\n", results)
	//
	//	if err != nil {
	//		fmt.Println(err.Error())
	//	}
}
