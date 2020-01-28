package main

import (
	"fmt"
	"github.com/GwangGwang/ganeungbot/pkg/translate"
)

func main() {
	fmt.Println("test")

	t, _ := translate.New("AIzaSyDimukIw-X7rckNqjuXLtl13Pu_Yc7WJZU")
	str, err := t.GetResponse("번역언어: ko")


	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(str)

	str, err = t.GetResponse("번역: blah blah")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(str)
}
