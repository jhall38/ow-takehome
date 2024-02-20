package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	OpenWeatherAPIVersion = "2.5"
	ColdThreshold         = 10 // Celsius
	ModerateThreshold     = 25
)

// WeatherResponse represents the relevant parts of the OpenWeather API response.
type WeatherResponse struct {
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

func (wr WeatherResponse) CombineWeatherDescriptions() string {
	if len(wr.Weather) == 0 {
		return "No specific weather data available"
	}

	var descriptions []string
	for _, w := range wr.Weather {
		description := w.Description
		if description == "" {
			description = w.Main
		}
		descriptions = append(descriptions, description)
	}
	return strings.Join(descriptions, ", ")
}

// Determines if the temperature is hot, cold, or moderate.
func (wr WeatherResponse) CategorizeTemperature() TemperatureCategory {
	temp := wr.Main.Temp
	switch {
	case temp <= ColdThreshold:
		return Cold
	case temp <= ModerateThreshold:
		return Moderate
	default:
		return Hot
	}
}

// TemperatureCategory describes the temperature as cold, moderate, or hot.
type TemperatureCategory string

const (
	Cold     TemperatureCategory = "cold"
	Moderate TemperatureCategory = "moderate"
	Hot      TemperatureCategory = "hot"
)

type OpenWeatherAPI interface {
	FetchWeather(apiKey, lat, lon string) (WeatherResponse, error)
}

type OpenWeatherAPIImpl struct{}

func (ow *OpenWeatherAPIImpl) FetchWeather(apiKey, lat, lon string) (WeatherResponse, error) {
	url := fmt.Sprintf("https://api.openweathermap.org/data/%s/weather?lat=%s&lon=%s&appid=%s&units=metric", OpenWeatherAPIVersion, lat, lon, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return WeatherResponse{}, fmt.Errorf("error fetching weather data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body) // assuming non-OK response is always relatively small
		if err != nil {
			log.Printf("error reading OpenWeather non-OK response body: %v", err)
		} else {
			log.Printf("OpenWeather non-OK Response body: %s", string(bodyBytes))
		}
		return WeatherResponse{}, fmt.Errorf("received non-OK status from OpenWeather API: %d", resp.StatusCode)
	}

	var weatherResponse WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResponse); err != nil {
		return WeatherResponse{}, fmt.Errorf("error decoding weather data: %v", err)
	}

	return weatherResponse, nil
}
