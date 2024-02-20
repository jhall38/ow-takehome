package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type WeatherDataResponse struct {
	Description         string              `json:"description"`
	TemperatureCategory TemperatureCategory `json:"temperature_category"`
}

// Handler for the /weather endpoint
func handleWeatherRequest(apiKey string, api OpenWeatherAPI, w http.ResponseWriter, r *http.Request) {
	lat := r.URL.Query().Get("lat")
	lon := r.URL.Query().Get("lon")

	if err := validateCoordinates(lat, lon); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	weatherResponse, err := api.FetchWeather(apiKey, lat, lon)
	if err != nil {
		log.Printf("Weather service failed: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	description := weatherResponse.CombineWeatherDescriptions()
	tempCategory := weatherResponse.CategorizeTemperature()

	response := WeatherDataResponse{
		Description:         description,
		TemperatureCategory: tempCategory,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// validateCoordinates checks if the provided latitude and longitude are valid.
func validateCoordinates(lat, lon string) error {

	if lat == "" || lon == "" {
		return fmt.Errorf("both latitude (\"lat\") and longitude (\"lon\") must be provided")
	}

	latFloat, err := strconv.ParseFloat(lat, 64)
	if err != nil || latFloat < -90 || latFloat > 90 {
		return fmt.Errorf("invalid latitude: %s", lat)
	}

	lonFloat, err := strconv.ParseFloat(lon, 64)
	if err != nil || lonFloat < -180 || lonFloat > 180 {
		return fmt.Errorf("invalid longitude: %s", lon)
	}

	return nil
}
