package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"os"
)

func main() {

	app := fiber.New()

	app.Use(logger.New())

	setupRoutes(app)

	if err := app.Listen(":8080"); err != nil {
		os.Exit(1)
	}
}

// setupRoutes registers all the application routes
func setupRoutes(app *fiber.App) {
	app.Get("/hw3/api", func(c *fiber.Ctx) error {
		if c.Query("queryStockPrice") != "" {
			return stockPriceHandler(c)
		} else if c.Query("queryAirportTemp") != "" {
			return airportTempHandler(c)
		} else if c.Query("queryEval") != "" {
			return evalHandler(c)
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid query parameter"})
		}
	})
}
