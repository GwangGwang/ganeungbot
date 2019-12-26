package main

import (
	"fmt"
	"github.com/GwangGwang/ganeungbot/pkg/geocoding"

)

func main() {
	fmt.Println("test")

	w, _ := geocoding.New("AIzaSyDmwAG0JTlIcZO0nFo6Tu7mQQwIZuhDiQQ")

	ll, str, _ := w.GetGPS("toronto")

	fmt.Printf("location: %+v", ll)
	fmt.Printf("formatted: %s", str)

	//	fmt.Printf("results: %+v\n", results)
	//
	//	if err != nil {
	//		fmt.Println(err.Error())
	//	}
}
