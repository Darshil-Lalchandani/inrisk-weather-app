package app

type FetchWeatherDataRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	//date in YYYY-MM-DD format
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type FetchWeatherDataResponse struct {
	Message string `json:"message"`
}

type ListWeatherFilesResponse struct {
	Files []string `json:"files"`
}
