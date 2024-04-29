package main

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	RapidAPIKey   string
	WeatherAPIKey string
	Port          string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	return &Config{
		RapidAPIKey:   os.Getenv("X_RAPIDAPI_KEY"),
		WeatherAPIKey: os.Getenv("WEATHERAPI_KEY"),
		Port:          os.Getenv("PORT"),
	}, nil
}
