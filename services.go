package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type StockAPIResponse struct {
	Status    string `json:"status"`
	RequestID string `json:"request_id"`
	Data      Data   `json:"data"`
}

type Data struct {
	Stock []StockItem `json:"stock"`
}

type StockItem struct {
	Symbol string  `json:"symbol"`
	Name   string  `json:"name"`
	Price  float64 `json:"price"`
}

func GetStockPrice(apiKey, query string) (*float64, error) {
	encodedQuery := url.QueryEscape(query)
	apiUrl := fmt.Sprintf("https://real-time-finance-data.p.rapidapi.com/search?query=%s&language=en", encodedQuery)

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-RapidAPI-Key", apiKey)
	req.Header.Add("X-RapidAPI-Host", "real-time-finance-data.p.rapidapi.com")

	log.Printf("Sending request to URL: %s with headers: %v", apiUrl, req.Header)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Println("Failed to close response body:", err)
		}
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response StockAPIResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if len(response.Data.Stock) == 0 {
		return nil, fmt.Errorf("no stock data found")
	}

	return &response.Data.Stock[0].Price, nil
}

type AirportInfo struct {
	ICAO      string `json:"icao"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type WeatherData struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func GetAirportTemp(apiKey, iata string) (*float64, error) {
	// Fetch airport info
	airportInfo, err := fetchAirportInfo(iata)
	if err != nil {
		return nil, err
	}
	log.Printf("Airport info: %+v", airportInfo)

	// Fetch weather info using the coordinates
	temp, err := fetchWeather(apiKey, airportInfo.Latitude, airportInfo.Longitude)
	if err != nil {
		return nil, err
	}

	return &temp, nil
}

func fetchAirportInfo(iata string) (*AirportInfo, error) {
	encodedQuery := url.QueryEscape(iata)
	apiUrl := fmt.Sprintf("https://www.airport-data.com/api/ap_info.json?iata=%s", encodedQuery)
	resp, err := http.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Failed to close response body:", err)
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		return nil, errors.New("failed to fetch airport info")
	}
	log.Printf("Response status code: %d", resp.StatusCode)

	var info AirportInfo
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &info)
	if err != nil {
		return nil, fmt.Errorf("airport not found")
	}

	if info.ICAO == "" || info.Latitude == "" || info.Longitude == "" {
		return nil, errors.New("failed to fetch airport info - data empty")
	}

	return &info, nil
}

func fetchWeather(apiKey, lat, lon string) (float64, error) {
	log.Printf("Fetching weather data for coordinates: %s, %s", lat, lon)
	apiUrl := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s,%s&aqi=no", apiKey, lat, lon)
	resp, err := http.Get(apiUrl)
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Failed to close response body:", err)
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		return 0, errors.New("failed to fetch weather data")
	}

	var weather WeatherData
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal(body, &weather)
	if err != nil {
		return 0, err
	}

	return weather.Current.TempC, nil
}
