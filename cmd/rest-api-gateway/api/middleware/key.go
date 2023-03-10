package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/krissukoco/go-microservices-marketplace/cmd/rest-api-gateway/config"
)

// x-api-key header middleware
// limit access to only clients with the correct api key, such as web app.
func CheckApiKeyHeader(c *fiber.Ctx) error {
	if c.Method() == "OPTIONS" || c.Path() == "/docs/*" {
		return c.Next()
	}
	apiKey := c.Get("x-api-key", "")
	// log.Println("x-api-key", apiKey, config.Api.ClientApiKey)
	if apiKey != config.Api.ClientApiKey {
		return c.SendStatus(401)
	}
	return c.Next()
}
