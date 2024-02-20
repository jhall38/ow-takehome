package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENWEATHER_API_KEY is not set in environment variables")
	}

	api := &OpenWeatherAPIImpl{}

	http.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		handleWeatherRequest(apiKey, api, w, r)
	})

	port := "8080"

	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
