package geocoding

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// Geocoding is the geocoding object
type Geocoding struct {
	GeocodingAPIKey string
	Operational     bool
	SavedData map[string]Location
}

// New initializes and returns a new geocoding Geocoding
func New(geocodingApiKey string) (Geocoding, error) {
	log.Println("Initializing geocoding pkg")

	g := Geocoding{}

	if len(geocodingApiKey) == 0 {
		log.Printf("WARN: geocoding API key not found")
		g.Operational = false
	} else {
		// TODO: test with toronto location
		g.Operational = true
		g.GeocodingAPIKey = geocodingApiKey
	}

	return g, nil
}

var userDefaultLocation = map[string]string{
	"younghaan": "toronto",
	"gwanggwang": "toronto",
	"appiejam": "toronto",
	"chanjook": "vancouver",
	"silversoup": "osaka",
}

func getDefaultUserLocation(username string) string {
	log.Printf("using default location for user '%s'", username)
	if loc, ok := userDefaultLocation[username]; ok {
		return loc
	}

	return ""
}

func (g *Geocoding) GetLocation(username string, loc string) (Location, string, error) {
	fmt.Printf("username: %s, location: %s", username, loc)

	var locationQuery string
	if len(loc) == 0 {
		locationQuery = getDefaultUserLocation(username)
	} else {
		locationQuery = loc
	}

	location, locStr, err := g.GetGPS(locationQuery)
	if err != nil {
		return Location{}, "", err
	}

	return location, locStr, nil
}

//GetGPS returns a location (lat/lng), formatted address, and error if any
func (g *Geocoding) GetGPS(location string) (Location, string, error) {
	handleError := func(err error) (Location, string, error) {
		return Location{}, "", err
	}

	url := fmt.Sprintf(geoCodingURL, url.QueryEscape(location), g.GeocodingAPIKey)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return handleError(fmt.Errorf("geocoding - error while sending request: %s", err))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var r Response
	err = json.Unmarshal(body, &r)
	if err != nil {
		return handleError(fmt.Errorf("geocoding - error while unmarshalling: %s", err))
	}
	fmt.Printf("resp: %+v\n\n", r)


	if r.Status != StatusOk {
		// TODO: some handling based on other status
		return handleError(fmt.Errorf("geocoding - non-ok status returned"))
	}

	if len(r.Results) < 1 {
		// does this even ever happen?
		return handleError(fmt.Errorf("geocoding - no results"))
	}

	// TODO: what if more than 1 locations returned?
	return r.Results[0].Geometry.Location, r.Results[0].FormattedAddress, nil
}
