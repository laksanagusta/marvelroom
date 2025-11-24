package http

import (
	"sandbox/internal/delivery/http/handler"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App, transactionHandler *handler.TransactionHandler, meetingHandler *handler.MeetingHandler, businessTripHandler *handler.BusinessTripHandler, assigneeHandler *handler.AssigneeHandler, businessTripTransactionHandler *handler.BusinessTripTransactionHandler) {
	api := app.Group("/api")

	api.Post("/upload", transactionHandler.UploadAndExtract)
	api.Post("/upload/detailed", transactionHandler.UploadAndExtractDetailed)
	api.Post("/report/excel", transactionHandler.GenerateRecapExcel)

	api.Post("/meetings", meetingHandler.CreateMeeting)

	api.Route("/v1/business-trips", func(r fiber.Router) {
		r.Post("/", businessTripHandler.CreateBusinessTrip)
		r.Get("/", businessTripHandler.ListBusinessTrips)
		r.Get("/:tripId", businessTripHandler.GetBusinessTrip)
		r.Put("/:tripId", businessTripHandler.UpdateBusinessTrip)
		r.Put("/:tripId/with-assignees", businessTripHandler.UpdateBusinessTripWithAssignees)
		r.Delete("/:tripId", businessTripHandler.DeleteBusinessTrip)

		r.Route("/:tripId/assignees", func(r fiber.Router) {
			r.Post("/", assigneeHandler.CreateAssignee)
			r.Get("/", assigneeHandler.ListAssignees)
			r.Get("/:assigneeId", assigneeHandler.GetAssignee)
			r.Put("/:assigneeId", assigneeHandler.UpdateAssignee)
			r.Delete("/:assigneeId", assigneeHandler.DeleteAssignee)

			r.Route("/:assigneeId/transactions", func(r fiber.Router) {
				r.Post("/", businessTripTransactionHandler.Create)
				r.Get("/", businessTripTransactionHandler.List)
				r.Put("/:transactionId", businessTripTransactionHandler.Update)
				r.Delete("/:transactionId", businessTripTransactionHandler.Delete)
			})
		})
	})

	// Legacy routes for backward compatibility
	if businessTripHandler != nil {
		businessTrips := api.Group("/business-trips")
		businessTrips.Get("/:id/summary", businessTripHandler.GetBusinessTripSummary)

		// Legacy assignee route
		businessTrips.Post("/:businessTripId/assignees", businessTripHandler.AddAssignee)

		// Legacy transaction routes
		assignees := api.Group("/assignees")
		assignees.Post("/:assigneeId/transactions", businessTripHandler.AddTransaction)
		assignees.Get("/:id/summary", businessTripHandler.GetAssigneeSummary)
	}

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
		})
	})
}
