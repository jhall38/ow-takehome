package main

import (
	"testing"
)

func TestCombineWeatherDescriptions(t *testing.T) {
	cases := []struct {
		name        string
		weatherResp WeatherResponse
		want        string
	}{
		{
			name: "one condition",
			weatherResp: WeatherResponse{
				Weather: []struct {
					Main        string `json:"main"`
					Description string `json:"description"`
				}{
					{Main: "Rain", Description: "light rain"},
				},
			},
			want: "light rain",
		},
		{
			name: "multiple conditions",
			weatherResp: WeatherResponse{
				Weather: []struct {
					Main        string `json:"main"`
					Description string `json:"description"`
				}{
					{Main: "Rain", Description: "light rain"},
					{Main: "Clouds", Description: "overcast clouds"},
				},
			},
			want: "light rain, overcast clouds",
		},
		{
			name: "empty weather array",
			weatherResp: WeatherResponse{
				Weather: []struct {
					Main        string `json:"main"`
					Description string `json:"description"`
				}{},
			},
			want: "No specific weather data available",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.weatherResp.CombineWeatherDescriptions()
			if got != c.want {
				t.Errorf("CombineWeatherDescriptions() = %q, want %q", got, c.want)
			}
		})
	}
}

func TestCategorizeTemperature(t *testing.T) {
	cases := []struct {
		name string
		temp float64
		want TemperatureCategory
	}{
		{"cold", 5, Cold},
		{"moderate", 15, Moderate},
		{"hot", 30, Hot},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			wr := WeatherResponse{
				Main: struct {
					Temp float64 `json:"temp"`
				}{Temp: c.temp},
			}
			got := wr.CategorizeTemperature()
			if got != c.want {
				t.Errorf("CategorizeTemperature() = %v, want %v", got, c.want)
			}
		})
	}
}
