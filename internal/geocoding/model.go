package geocoding

type GeocodingStatus string
const (
	StatusOk GeocodingStatus = "OK"
	StatusZeroResults = "ZERO_RESULTS" // indicates that the geocode was successful but returned no results. This may occur if the geocoder was passed a non-existent address.
	StatusOverDailyLimit = "OVER_DAILY_LIMIT" // api key missing / billing not enabled / self-imposed usage cap
	StatusOverQueryLimit = "OVER_QUERY_LIMIT" // indicates that you are over your quota.
	StatusRequestDenied = "REQUEST_DENIED" // indicates that your request was denied.
	StatusInvalidRequest = "INVALID_REQUEST" // generally indicates that the query (address, components or latlng) is missing.
	StatusUnknownError = "UNKNOWN_ERROR"
)

const geoCodingURL = "https://maps.googleapis.com/maps/api/geocode/json?address=%s&key=%s"

//type Type string

type Response struct {
	Results []Result `json:"results"`
	Status GeocodingStatus `json:"status"`
}

type Result struct {
	PlaceID string `json:"place_id"`
	Types []string `json:"types"`
	AddressComponents []AddressComponent `json:"address_components"`
	FormattedAddress string `json:"formatted_address"`
	Geometry Geometry `json:"geometry"`
}

type AddressComponent struct {
	//Types []Type `json:"types"`
	LongName string `json:"long_name"`
	ShortName string `json:"short_name"`
}

type Geometry struct {
	Location Location `json:"location"`
	//LocationType string `json:"location_type"`
	//Viewport Bounds `json:"viewport"`
	//Bounds Bounds `json:"bounds"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

//type Bounds struct {
//	Northeast Location `json:"northeast"`
//	Southwest Location `json:"southwest"`
//}


