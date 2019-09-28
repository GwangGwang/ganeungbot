package weather

// Weather forecasting via DarkSky API

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	darksky "github.com/twpayne/go-darksky"
)

type API struct {
	WeatherAPIKey   string
	GeoCodingAPIKey string
}

const GeoCodingURL = "http://www.mapquestapi.com/geocoding/v1/address?key=%s&location=%s"

type GeoCodingResults struct {
	Results []Result `"json:results"`
}

type Result struct {
	Locations []Location `"json:locations"`
}

type Location struct {
	LatLng LatLng `"json:latlng"`
}

type LatLng struct {
	Lat float64 `"json:lat"`
	Lng float64 `"json:lng"`
}

func (w *API) Start() {

	geocodingURL := fmt.Sprintf(GeoCodingURL, w.GeoCodingAPIKey, "Toronto")
	resp, err := http.Get(geocodingURL)
	respBody, err := ioutil.ReadAll(resp.Body)

	var gr GeoCodingResults
	err = json.Unmarshal(respBody, &gr)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%+v", gr)
	}

	c, err := darksky.NewClient(
		darksky.WithKey(w.WeatherAPIKey),
	)
	ctx := context.Background()
	forecast, err := c.Forecast(ctx, 42.3601, -71.0589, nil, &darksky.ForecastOptions{
		Units: darksky.UnitsSI,
	})

	if err == nil {
	} else {
		fmt.Println(err)
	}
	// The forecast varies from day to day. Print something stable.
	fmt.Println(forecast.Timezone)
}
