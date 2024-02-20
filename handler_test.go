package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockOpenWeatherAPI struct {
	FetchWeatherFunc func(apiKey, lat, lon string) (WeatherResponse, error)
}

func (m *MockOpenWeatherAPI) FetchWeather(apiKey, lat, lon string) (WeatherResponse, error) {
	return m.FetchWeatherFunc(apiKey, lat, lon)
}

func TestHandleWeatherRequest(t *testing.T) {
	mockAPI := &MockOpenWeatherAPI{
		FetchWeatherFunc: func(apiKey, lat, lon string) (WeatherResponse, error) {
			return WeatherResponse{
				Weather: []struct {
					Main        string `json:"main"`
					Description string `json:"description"`
				}{
					{Main: "Clouds", Description: "overcast clouds"},
				},
				Main: struct {
					Temp float64 `json:"temp"`
				}{Temp: 15},
			}, nil
		},
	}

	req, err := http.NewRequest("GET", "/weather?lat=35&lon=139", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleWeatherRequest("foo", mockAPI, w, r)
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := WeatherDataResponse{
		Description:         "overcast clouds",
		TemperatureCategory: Moderate,
	}
	var actual WeatherDataResponse
	if err := json.NewDecoder(rr.Body).Decode(&actual); err != nil {
		t.Fatal("Failed to decode response body")
	}

	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}

func TestValidateCoordinates(t *testing.T) {
	cases := []struct {
		name      string
		lat, lon  string
		wantError bool
	}{
		{"valid coordinates", "40.712776", "-74.005974", false},
		{"latitude high", "91.0", "-74.005974", true},
		{"latitude low", "-91.0", "-74.005974", true},
		{"longitude high", "40.712776", "181.0", true},
		{"longitude low", "40.712776", "-181.0", true},
		{"non numeric latitude", "nonnumber", "-74.005974", true},
		{"non numeric longitude", "40.712776", "nonnumber", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateCoordinates(tc.lat, tc.lon)
			if tc.wantError && err == nil {
				t.Errorf("validateCoordinates(%s, %s) expected an error, but got none", tc.lat, tc.lon)
			} else if !tc.wantError && err != nil {
				t.Errorf("validateCoordinates(%s, %s) unexpected error: %v", tc.lat, tc.lon, err)
			}
		})
	}
}
