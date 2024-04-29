package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"os"
)

func main() {

	config, err := LoadConfig()
	if err != nil {
		log.Print("Failed to load configuration")
		os.Exit(1)
	}
	if config.RapidAPIKey == "" || config.WeatherAPIKey == "" {
		log.Printf("Missing configuration values: RapidAPIKey=%s, WeatherAPIKey=%s", config.RapidAPIKey, config.WeatherAPIKey)
		os.Exit(1)
	}
	if config.Port == "" {
		config.Port = "16443"
	}

	app := fiber.New()
	app.Use(logger.New())
	setupRoutes(app, config)

	if err := app.Listen("0.0.0.0:" + config.Port); err != nil {
		os.Exit(1)
	}
}

func setupRoutes(app *fiber.App, config *Config) {
	app.Get("/hw3/api", func(c *fiber.Ctx) error {
		if c.Query("queryStockPrice") != "" {
			return stockPriceHandler(c, config.RapidAPIKey)
		} else if c.Query("queryAirportTemp") != "" {
			return airportTempHandler(c, config.WeatherAPIKey)
		} else if c.Query("queryEval") != "" {
			return evalHandler(c)
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid query parameter"})
		}
	})
}
