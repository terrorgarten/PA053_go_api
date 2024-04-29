package main

import (
	"github.com/Knetic/govaluate"
	"github.com/gofiber/fiber/v2"
	"log"
)

func stockPriceHandler(c *fiber.Ctx, stockAPIKey string) error {

	query := c.Query("queryStockPrice")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "query parameter is required"})
	}

	if len(query) < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "query parameter must be at least 1 character"})
	}
	if len(query) > 4 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "query parameter must be at most 5 characters"})
	}

	log.Printf("Query received for stock price: %s", query)

	stockPrice, err := GetStockPrice(stockAPIKey, query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch stock price", "details": err.Error()})
	}

	return c.JSON(stockPrice)
}

func airportTempHandler(c *fiber.Ctx, weatherAPIKey string) error {
	query := c.Query("queryAirportTemp")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "query parameter is required"})
	}

	log.Printf("Query received for airport temperature: %s", query)

	temperature, err := GetAirportTemp(weatherAPIKey, query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch airport temperature", "details": err.Error()})
	}

	return c.JSON(temperature)
}

func evalHandler(c *fiber.Ctx) error {
	expression, err := govaluate.NewEvaluableExpression(c.Query("queryEval"))
	if err != nil {
		log.Printf("Failed to parse expression: %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid expression"})
	}

	result, err := expression.Evaluate(nil)
	if err != nil {
		log.Printf("Failed to evaluate expression: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to evaluate expression"})
	}

	resultFloat, ok := result.(float64)
	if !ok {
		log.Printf("Failed to convert result to float64: %v", result)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to convert result to float64"})
	}

	return c.Status(fiber.StatusOK).JSON(resultFloat)
}
