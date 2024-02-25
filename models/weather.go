package models

type WeatherData struct {
	Current struct {
			Cloud            int     `json:"cloud"`
			FeelsLikeC       float64 `json:"feelslike_c"`
			FeelsLikeF       float64 `json:"feelslike_f"`
			GustKph          float64 `json:"gust_kph"`
			GustMph          float64 `json:"gust_mph"`
			Humidity         int     `json:"humidity"`
			IsDay            int     `json:"is_day"`
			LastUpdated      string  `json:"last_updated"`
			LastUpdatedEpoch int     `json:"last_updated_epoch"`
			PrecipIn         float64 `json:"precip_in"`
			PrecipMm         float64 `json:"precip_mm"`
			PressureIn       float64 `json:"pressure_in"`
			PressureMb       float64 `json:"pressure_mb"`
			TempC            float64 `json:"temp_c"`
			TempF            float64 `json:"temp_f"`
			UV               float64 `json:"uv"`
			VisKm            float64 `json:"vis_km"`
			VisMiles         float64 `json:"vis_miles"`
			WindDegree       int     `json:"wind_degree"`
			WindDir          string  `json:"wind_dir"`
			WindKph          float64 `json:"wind_kph"`
			WindMph          float64 `json:"wind_mph"`
			Condition        struct {
					Code         int     `json:"code"`
					Icon         string  `json:"icon"`
					Text         string  `json:"text"`
			}                        `json:"condition"`
	}                            `json:"current"`
	Location struct {
			Country          string  `json:"country"`
			Lat              float64 `json:"lat"`
			Localtime        string  `json:"localtime"`
			LocaltimeEpoch   int64   `json:"localtime_epoch"`
			Lon              float64 `json:"lon"`
			Name             string  `json:"name"`
			Region           string  `json:"region"`
			TzID             string  `json:"tz_id"`
	}                            `json:"location"`
}