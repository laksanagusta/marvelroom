package main

import (
	"fmt"
	"log"

	"sandbox/config"
	httpRouter "sandbox/internal/delivery/http"
	"sandbox/internal/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize dependency injection container
	container := config.NewContainer(cfg)

	// Setup Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler,
	})

	// Setup middleware
	app.Use(middleware.ConfigureLogger())
	app.Use(middleware.ConfigureRecovery())
	app.Use(middleware.ConfigureCORS(cfg.CORS.AllowOrigins))

	// Setup routes with all handlers
	httpRouter.SetupRoutes(app, container.TransactionHandler, container.MeetingHandler, container.BusinessTripHandler, container.AssigneeHandler, container.BusinessTripTransactionHandler)

	// Start server
	fmt.Printf("🚀 Server running on port %s\n", cfg.Server.Port)
	fmt.Printf("📝 Environment: %s\n", getEnvironment())

	if err := app.Listen(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// customErrorHandler handles errors globally
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
		"code":  code,
	})
}

func getEnvironment() string {
	return "development"
}
