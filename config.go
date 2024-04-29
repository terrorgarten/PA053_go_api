package main

import (
	"os"
)

type Config struct {
	RapidAPIKey   string
	WeatherAPIKey string
	Port          string
}

func LoadConfig() (*Config, error) {
	return &Config{
		RapidAPIKey:   os.Getenv("X_RAPIDAPI_KEY"),
		WeatherAPIKey: os.Getenv("WEATHERAPI_KEY"),
		Port:          os.Getenv("PORT"),
	}, nil
}
