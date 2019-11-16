package weather

// A Forecast is a forecast.
//type Forecast struct {
//	Alerts    []*Alert  `json:"alerts"`
//	Currently *DataSet  `json:"currently"`
//	Daily     *DataSet  `json:"daily"`
//	Hourly    *DataSet  `json:"hourly"`
//	Minutely  *Minutely `json:"minutely"`
//	Latitude  float64   `json:"latitude"`
//	Longitude float64   `json:"longitude"`
//	Flags     *Flags    `json:"flags"`
//	Offset    float64   `json:"offset"`
//	Timezone  string    `json:"timezone"`
//}
//
//type DataSet struct {
//	Data    []Data `json:"data"`
//	Icon    Icon   `json:"icon"`
//	Summary string `json:"summary"`
//}
//
//type Data struct {
//	ApparentTemperature         Measurement `json:"apparentTemperature,omitempty"`
//	ApparentTemperatureHigh     Measurement `json:"apparentTemperatureHigh,omitempty"`
//	ApparentTemperatureHighTime Timestamp   `json:"apparentTemperatureHighTime,omitempty"`
//	ApparentTemperatureLow      Measurement `json:"apparentTemperatureLow,omitempty"`
//	ApparentTemperatureLowTime  Timestamp   `json:"apparentTemperatureLowTime,omitempty"`
//	ApparentTemperatureMax      Measurement `json:"apparentTemperatureMax,omitempty"`
//	ApparentTemperatureMaxTime  Timestamp   `json:"apparentTemperatureMaxTime,omitempty"`
//	ApparentTemperatureMin      Measurement `json:"apparentTemperatureMin,omitempty"`
//	ApparentTemperatureMinTime  Timestamp   `json:"apparentTemperatureMinTime,omitempty"`
//	CloudCover                  Measurement `json:"cloudCover,omitempty"`
//	DewPoint                    Measurement `json:"dewPoint,omitempty"`
//	Humidity                    Measurement `json:"humidity,omitempty"`
//	Icon                        string      `json:"icon,omitempty"`
//	MoonPhase                   Measurement `json:"moonPhase,omitempty"`
//	NearestStormBearing         Measurement `json:"nearestStormBearing,omitempty"`
//	NearestStormDistance        Measurement `json:"nearestStormDistance,omitempty"`
//	Ozone                       Measurement `json:"ozone,omitempty"`
//	PrecipAccumulation          Measurement `json:"precipAccumulation,omitempty"`
//	PrecipIntensity             Measurement `json:"precipIntensity,omitempty"`
//	PrecipIntensityError        Measurement `json:"precipIntensityError,omitempty"`
//	PrecipIntensityMax          Measurement `json:"precipIntensityMax,omitempty"`
//	PrecipIntensityMaxTime      Timestamp   `json:"precipIntensityMaxTime,omitempty"`
//	PrecipProbability           Measurement `json:"precipProbability,omitempty"`
//	PrecipType                  string      `json:"precipType,omitempty"`
//	Pressure                    Measurement `json:"pressure,omitempty"`
//	Summary                     string      `json:"summary,omitempty"`
//	SunriseTime                 Timestamp   `json:"sunriseTime,omitempty"`
//	SunsetTime                  Timestamp   `json:"sunsetTime,omitempty"`
//	Temperature                 Measurement `json:"temperature,omitempty"`
//	TemperatureHigh             Measurement `json:"temperatureHigh,omitempty"`
//	TemperatureHighTime         Timestamp   `json:"temperatureHighTime,omitempty"`
//	TemperatureLow              Measurement `json:"temperatureLow,omitempty"`
//	TemperatureLowTime          Timestamp   `json:"temperatureLowTime,omitempty"`
//	TemperatureMax              Measurement `json:"temperatureMax,omitempty"`
//	TemperatureMaxTime          Timestamp   `json:"temperatureMaxTime,omitempty"`
//	TemperatureMin              Measurement `json:"temperatureMin,omitempty"`
//	TemperatureMinTime          Timestamp   `json:"temperatureMinTime,omitempty"`
//	Time                        Timestamp   `json:"time,omitempty"`
//	UvIndex                     int64       `json:"uvIndex,omitempty"`
//	UvIndexTime                 Timestamp   `json:"uvIndexTime,omitempty"`
//	Visibility                  Measurement `json:"visibility,omitempty"`
//	WindBearing                 Measurement `json:"windBearing,omitempty"`
//	WindGust                    Measurement `json:"windGust,omitempty"`
//	WindGustTime                Timestamp   `json:"windGustTime,omitempty"`
//	WindSpeed                   Measurement `json:"windSpeed,omitempty"`
//}
//
//// Flags are forecast flags.
//type Flags struct {
//	DarkSkyUnavailable interface{} `json:"darksky-unavailable"`
//	NearestStation     float64     `json:"nearest-station"`
//	Sources            []string    `json:"sources"`
//	Units              Units       `json:"units"`
//}
//
//type Alert struct {
//	Description string    `json:"description"`
//	Expires     TimeStamp `json:"expires"`
//	Regions     []string  `json:"regions"`
//	Severity    Severity  `json:"severity"`
//	Time        TimeStamp `json:"time"`
//	Title       string    `json:"title"`
//	URI         string    `json:"uri"`
//}

//func (w *Instance) GetForecast(lat float64, lng float64) (darksky.Forecast, error) {
//
//	fmt.Println(lat)
//	fmt.Println(lng)
//	c, err := darksky.NewClient(
//		darksky.WithKey(w.WeatherAPIKey),
//	)
//	ctx := context.Background()
//	forecast, err := c.Forecast(ctx, lat, lng, nil, &darksky.ForecastOptions{
//		Units: darksky.UnitsSI,
//	})
//
//	if err != nil {
//		return darksky.Forecast{}, err
//	}
//	// The forecast varies from day to day. Print something stable.
//
//	return *forecast, nil
//}
