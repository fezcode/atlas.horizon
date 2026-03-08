package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Location struct {
	City      string  `json:"city"`
	Country   string  `json:"country"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	Timezone  string  `json:"timezone"`
}

type WeatherData struct {
	Current struct {
		Temperature   float64 `json:"temperature_2m"`
		Humidity      float64 `json:"relative_humidity_2m"`
		ApparentTemp  float64 `json:"apparent_temperature"`
		IsDay         int     `json:"is_day"`
		WeatherCode   int     `json:"weather_code"`
		WindSpeed     float64 `json:"wind_speed_10m"`
		WindDirection int     `json:"wind_direction_10m"`
	} `json:"current"`
	Daily struct {
		Time           []string  `json:"time"`
		WeatherCode    []int     `json:"weather_code"`
		TempMax        []float64 `json:"temperature_2m_max"`
		TempMin        []float64 `json:"temperature_2m_min"`
		UVIndex        []float64 `json:"uv_index_max"`
	} `json:"daily"`
}

func GetLocation() (*Location, error) {
	resp, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var loc Location
	if err := json.NewDecoder(resp.Body).Decode(&loc); err != nil {
		return nil, err
	}
	return &loc, nil
}

func GetWeather(lat, lon float64, timezone string) (*WeatherData, error) {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m,relative_humidity_2m,apparent_temperature,is_day,weather_code,wind_speed_10m,wind_direction_10m&daily=weather_code,temperature_2m_max,temperature_2m_min,uv_index_max&timezone=%s", lat, lon, timezone)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}

func GetWeatherCondition(code int) string {
	switch code {
	case 0: return "Clear Sky"
	case 1, 2, 3: return "Partly Cloudy"
	case 45, 48: return "Foggy"
	case 51, 53, 55: return "Drizzle"
	case 61, 63, 65: return "Rain"
	case 71, 73, 75: return "Snow"
	case 95, 96, 99: return "Thunderstorm"
	default: return "Unknown"
	}
}

func GetWeatherIcon(code int, isDay bool) string {
	switch code {
	case 0: 
		if isDay { return "☀️" }
		return "🌙"
	case 1, 2, 3: 
		if isDay { return "⛅" }
		return "☁️"
	case 45, 48: return "🌫️"
	case 51, 53, 55: return "🌦️"
	case 61, 63, 65: return "🌧️"
	case 71, 73, 75: return "❄️"
	case 95, 96, 99: return "⛈️"
	default: return "❓"
	}
}

func Geocode(query string) (*Location, error) {
	url := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1&language=en&format=json", query)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Results []struct {
			Name      string  `json:"name"`
			Country   string  `json:"country"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			Timezone  string  `json:"timezone"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Results) == 0 {
		return nil, fmt.Errorf("location not found")
	}

	r := result.Results[0]
	return &Location{
		City:     r.Name,
		Country:  r.Country,
		Lat:      r.Latitude,
		Lon:      r.Longitude,
		Timezone: r.Timezone,
	}, nil
}
