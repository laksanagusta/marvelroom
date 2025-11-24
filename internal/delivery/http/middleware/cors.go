package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func ConfigureCORS(allowOrigins string) fiber.Handler {
	log.Println(allowOrigins)
	if allowOrigins == "" {
		allowOrigins = "https://marvcore.com,https://www.marvcore.com,http://localhost:3000"
	}

	return cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	})
}
