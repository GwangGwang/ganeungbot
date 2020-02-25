package weather

// Weather forecasting via MapQuest Geocoding API and DarkSky API

import (
	"context"
	"fmt"
	"github.com/GwangGwang/ganeungbot/internal/geocoding"
	"log"
	"time"

	darksky "github.com/twpayne/go-darksky"
)

// Weather is the weather forecast object
type Weather struct {
	WeatherAPIKey   string
	Geocoding       geocoding.Geocoding
	Operational     bool
}

// New initializes and returns a new weather pkg Weather
func New(weatherApiKey string, geocoding geocoding.Geocoding) (Weather, error) {
	log.Println("Initializing weather pkg")

	w := Weather{}

	if len(weatherApiKey) == 0 {
		log.Printf("WARN: weather API key not found")
		w.Operational = false
	} else {
		w.Operational = true
	}

	w.Geocoding = geocoding
	w.WeatherAPIKey = weatherApiKey

	return w, nil
}

// GetResponse is the main outward facing function to generate weather response
func (w *Weather) GetResponse(username string, txt string) (string, error) {
	// parse out time/location keywords and process any time offsets
	parseResult, err := parse(txt)
	if err != nil {
		return "", err
	}

	queryLocation := parseResult.Location
	location, locStr, err := w.Geocoding.GetLocation(username, queryLocation)
	if err != nil {
		fmt.Printf("error in geocoding: %s", err)
		return "", fmt.Errorf("error in geocoding api")
	}

	forecast, err := w.GetForecast(location.Lat, location.Lng)
	if err != nil {
		fmt.Printf("error in forecast: %s", err)
		return "", fmt.Errorf("error in forecast api")
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

func (w *Weather) GetForecast(lat float64, lng float64) (darksky.Forecast, error) {

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

func (w *Weather) BuildResponse(f darksky.Forecast, pr parseResult) (string, error) {

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
