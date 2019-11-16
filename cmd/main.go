package main

import "github.com/GwangGwang/ganeungbot/pkg/weather"

func main() {

	w := weather.APIInfo{
		WeatherAPIKey:   "d0c0819bee7bc85b2cbc33106dad8125",
		GeoCodingAPIKey: "oSFvEcUqRbvWl7C6VB8DcMVNZrGphOMw",
	}

	weather.Start(w)
}
