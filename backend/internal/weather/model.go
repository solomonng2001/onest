package weather

// APIResponse represents the fields we use from the data.gov.sg response.
type APIResponse struct {
	Code     int     `json:"code"`
	Data     APIData `json:"data"`
	ErrorMsg string  `json:"errorMsg"`
}

type APIData struct {
	Items           []ForecastItem `json:"items"`
	PaginationToken string         `json:"paginationToken"`
}

type ForecastItem struct {
	ValidPeriod ValidPeriod    `json:"valid_period"`
	Forecasts   []AreaForecast `json:"forecasts"`
}

type ValidPeriod struct {
	Text string `json:"text"`
}

type AreaForecast struct {
	Area     string `json:"area"`
	Forecast string `json:"forecast"`
}

// LocationForecast represents one row in the frontend table.
type LocationForecast struct {
	Location string `json:"location"`
	Forecast string `json:"forecast"`
}

// ForecastResponse represents the complete response sent to the frontend.
type ForecastResponse struct {
	ValidPeriod string             `json:"validPeriod"`
	Forecasts   []LocationForecast `json:"forecasts"`
}
