package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const geoCodingURL = "http://www.mapquestapi.com/geocoding/v1/address?key=%s&location=%s"

type GeoCodingResults struct {
	Results []Result `"json:results"`
}

type Result struct {
	Locations []Location `"json:locations"`
}

type Location struct {
	LatLng     LatLng `"json:latlng"`
	AdminArea6 string `"json:adminArea6"` // Neighborhood
	AdminArea5 string `"json:adminArea5"` // City
	AdminArea4 string `"json:adminArea4"` // County
	AdminArea3 string `"json:adminArea3"` // State
	AdminArea1 string `"json:adminArea1"` // Country
}

type LatLng struct {
	Lat float64 `"json:lat"`
	Lng float64 `"json:lng"`
}

func (w *Instance) GetGPS(location string) (LatLng, string, error) {
	url := fmt.Sprintf(geoCodingURL, w.GeocodingAPIKey, location)
	resp, err := http.Get(url)
	if err != nil {
		return LatLng{}, "", fmt.Errorf("error while sending request: %s", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var r GeoCodingResults
	err = json.Unmarshal(body, &r)
	if err != nil {
		return LatLng{}, "", fmt.Errorf("geocoding - error while unmarshalling: %s", err)
	}

	// TODO: some validation stuff
	if len(r.Results[0].Locations) < 1 {
		return LatLng{}, "", fmt.Errorf("no results from geocoding api")
	}

	// city, country
	locStr := fmt.Sprintf("%s, %s", r.Results[0].Locations[0].AdminArea5, r.Results[0].Locations[0].AdminArea1)

	// TODO: what if more than 1 locations returned?
	return r.Results[0].Locations[0].LatLng, locStr, nil
}
