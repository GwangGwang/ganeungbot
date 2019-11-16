package weather

// Weather forecasting via MapQuest Geocoding API and DarkSky API

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/GwangGwang/ganeungbot/internal/pkg/config"
	darksky "github.com/twpayne/go-darksky"
)

// Instance is the weather forecast object
type Instance struct {
	WeatherAPIKey   string
	GeocodingAPIKey string
	Operational     bool
	LocationLatLngMap map[string]LatLng
}

const weatherAPIKey = "weatherAPIKey"
const geocodingAPIKey = "geocodingAPIKey"

// New initializes and returns a new weather pkg instance
func New() (Instance, error) {
	log.Println("Initializing weather pkg")

	w := Instance{}

	weatherAPIKey, err := config.Get(weatherAPIKey)
	if err != nil {
		log.Printf("WARN: weather API key not found: %s", err.Error())
		w.Operational = false
	}
	geocodingAPIKey, err := config.Get(geocodingAPIKey)
	if err != nil {
		log.Printf("WARN: geocoding API key not found: %s", err.Error())
		w.Operational = false
	}

	w.Operational = true
	w.WeatherAPIKey = weatherAPIKey
	w.GeocodingAPIKey = geocodingAPIKey

	return w, nil
}

func (w *Instance) getLocation(loc string) (LatLng, string, error) {
	var locationQuery string

	// see if there's a pre-set location name
	// TODO: save geo coding queries somewhere and retrieve
	defaultLoc := getSavedLocation(loc)
	if len(defaultLoc) > 0 {
		locationQuery = defaultLoc
	} else {
		locationQuery = loc
	}

	latlng, locStr, err := w.GetGPS(locationQuery)
	if err != nil {
		return LatLng{}, "", err
	}

	return latlng, locStr, nil
}

var userDefaultLocation map[string]string = map[string]string{
	"younghaan": "toronto",
	"gwanggwang": "toronto",
	"appiejam": "toronto",
	"chanjook": "vancouver",
	"silversoup": "osaka",
}

func (w *Instance) getDefaultUserLocation(username string) string {
	log.Printf("using default location for user '%s'", username)
	if loc, ok := userDefaultLocation[username]; ok {
		return loc
	}

	return ""
}

// GetResponse is the main outward facing function to generate weather response
func (w *Instance) GetResponse(username string, txt string) (string, error) {
	// parse out time/location keywords and process any time offsets
	parseResult, err := parse(txt)
	if err != nil {
		return "", err
	}

	queryLocation := parseResult.Location
	if len(queryLocation) == 0 {
		// no location supplied in weather query; search for default location for the user
		queryLocation = w.getDefaultUserLocation(username)
	}

	latlng, locStr, err := w.getLocation(queryLocation)
	if err != nil {
		return "", err
	}

	forecast, err := w.GetForecast(latlng.Lat, latlng.Lng)
	if err != nil {
		return "", err
	}

	fmt.Printf("currently:\n %+v\n", forecast.Currently)
	fmt.Printf("daily:\n %+v\n", forecast.Daily)
	fmt.Printf("hourly:\n %+v\n", forecast.Hourly)

	parseResult.Location = locStr
	resp, err := w.BuildResponse(forecast, parseResult)
	if err != nil {
		return "", err
	}

	return resp, nil
}

func (w *Instance) GetForecast(lat float64, lng float64) (darksky.Forecast, error) {

	fmt.Println(lat)
	fmt.Println(lng)
	c, err := darksky.NewClient(
		darksky.WithKey(w.WeatherAPIKey),
	)
	ctx := context.Background()
	forecast, err := c.Forecast(ctx, lat, lng, nil, &darksky.ForecastOptions{
		Units: darksky.UnitsSI,
	})

	if err != nil {
		return darksky.Forecast{}, err
	}
	// The forecast varies from day to day. Print something stable.

	return *forecast, nil
}

func (w *Instance) BuildResponse(f darksky.Forecast, pr parseResult) (string, error) {

	var resp string

	// Get timezone
	timezone, err := time.LoadLocation(f.Timezone)
	if err != nil {
		return "", fmt.Errorf("error from loading timezone from darksky: %s", err)
	}

	// Get time object for query time
	// 0am for that day + offset
	curYear, curMonth, curDay := time.Now().Date()
	queryTime := time.Date(curYear, curMonth, curDay, 0, 0, 0, 0, timezone)
	dur, _ := time.ParseDuration(fmt.Sprintf("%ds", pr.TimeInfo.Offset))
	queryTime = queryTime.Add(dur)

	_, month, date := queryTime.Date()

	var summary string
	var temp, apparentTemp, humidity, pop float64

	switch pr.TimeInfo.Category {
	case now:
		resp = fmt.Sprintf("<현재 %s 날씨>\n", pr.Location)

		curForecast := f.Currently
		temp = (curForecast.Temperature)
		apparentTemp = (curForecast.ApparentTemperature)
		humidity = curForecast.Humidity * 100
		pop = curForecast.PrecipProbability * 100 // percentage
		summary = curForecast.Summary

		if summary != "" {
			resp += summary + "\n"
		}
		resp += fmt.Sprintf("Temperature: %.1f°C (feels like %.1f°C)\n", temp, apparentTemp) +
			fmt.Sprintf("Humidity: %.0f%%\n", humidity) +
			fmt.Sprintf("POP: %.0f%%", pop)
	case day:
		//title = fmt.Sprintf("<%d/%d/%d %s 날씨>", year, month, date, pr.Location)
		resp = fmt.Sprintf("<%d/%d %s 날씨>\n", month, date, pr.Location)

		for _, data := range f.Daily.Data {
			if queryTime.Equal(data.Time.Time) {
				fmt.Printf("match!")
				humidity = data.Humidity * 100
				pop = data.PrecipProbability * 100 // percentage

				if data.Summary != "" {
					resp += data.Summary + "\n"
				}
				resp += fmt.Sprintf("Temperature: %.1f/%.1f°C\n", (data.TemperatureHigh), (data.TemperatureLow)) +
					fmt.Sprintf("   (Feels like %.1f/%.1f°C)\n", (data.ApparentTemperatureHigh), (data.ApparentTemperatureLow)) +
					fmt.Sprintf("Humidity: %.0f%%\n", humidity) +
					fmt.Sprintf("POP: %.0f%%", pop)
			}
		}
	case hour:
		hour := queryTime.Hour()

		var hourTxt string
		if hour == 12 {
			hourTxt = "12pm"
		} else if hour >= 12 {
			hourTxt = fmt.Sprintf("%dpm", hour-12)
		} else {
			hourTxt = fmt.Sprintf("%dam", hour)
		}

		//title = fmt.Sprintf("%d/%d/%d %s %s 날씨", year, month, date, hourTxt, pr.Location)
		resp = fmt.Sprintf("%d/%d %s %s 날씨\n", month, date, hourTxt, pr.Location)

		for _, data := range f.Hourly.Data {
			//fmt.Printf("q time: %s\nd time: %s\n\n", queryTime.String(), data.Time.Time.String())
			if queryTime.Equal(data.Time.Time) {
				summary := data.Summary
				if summary != "" {
					resp += summary + "\n"
				}
				temp := data.Temperature
				apparentTemp := data.ApparentTemperature
				humidity := data.Humidity * 100
				pop := data.PrecipProbability * 100 // percentage

				resp += fmt.Sprintf("Temperature: %.1f°C (Feels like %.1f°C)\n", temp, apparentTemp) +
					fmt.Sprintf("Humidity: %.0f%%\n", humidity) +
					fmt.Sprintf("POP: %.0f%%", pop)
			}
		}
	}

	return resp, nil
}

func fToC(f float64) float64 {
	return (f - 32) * 5 / 9
}
