package http

import (
	"sandbox/internal/delivery/http/handler"
	"sandbox/internal/delivery/http/middleware"
	deskHandler "sandbox/internal/delivery/http/handler/desk"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App, transactionHandler *handler.TransactionHandler, meetingHandler *handler.MeetingHandler, businessTripHandler *handler.BusinessTripHandler, assigneeHandler *handler.AssigneeHandler, businessTripTransactionHandler *handler.BusinessTripTransactionHandler, workPaperItemHandler *deskHandler.WorkPaperItemHandler, workPaperHandler *deskHandler.WorkPaperHandler, vaccineHandler *handler.VaccineHandler, signatureHandler *handler.WorkPaperSignatureHandler, cryptoHandler *handler.CryptoHandler) {
	api := app.Group("/api")

	api.Post("/upload", middleware.AuthMiddleware(), transactionHandler.UploadAndExtract)
	api.Post("/upload/detailed", middleware.AuthMiddleware(), transactionHandler.UploadAndExtractDetailed)
	api.Post("/report/excel", middleware.AuthMiddleware(), transactionHandler.GenerateRecapExcel)

	api.Post("/meetings", middleware.AuthMiddleware(), meetingHandler.CreateMeeting)

	api.Route("/v1/business-trips", func(r fiber.Router) {
		r.Use(middleware.AuthMiddleware()) // Apply auth middleware to all business trips routes
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

	// Desk module routes
	api.Route("/v1/desk", func(r fiber.Router) {
		r.Use(middleware.AuthMiddleware()) // Apply auth middleware to all desk routes
		// Work Paper Item routes (new)
		r.Route("/work-paper-items", func(r fiber.Router) {
			r.Post("/", workPaperItemHandler.CreateWorkPaperItem)
			r.Get("/", workPaperItemHandler.ListWorkPaperItems)
		})

		// Work Paper routes (new)
		r.Route("/work-papers", func(r fiber.Router) {
			r.Post("/", workPaperHandler.CreateWorkPaper)
			r.Get("/", workPaperHandler.ListWorkPapers)
			r.Get("/status-transitions", workPaperHandler.GetStatusTransitions)
			r.Get("/:id", workPaperHandler.GetWorkPaperByID)
			r.Put("/:id/status", workPaperHandler.UpdateWorkPaperStatus)
		})

		// Work Paper Note routes (new)
		r.Post("/work-paper-notes/check", workPaperHandler.CheckWorkPaperNote)
		r.Put("/work-paper-notes/:id", workPaperHandler.UpdateWorkPaperNote)

		// Work Paper Signature routes
		r.Route("/work-paper-signatures", func(r fiber.Router) {
			r.Post("/", signatureHandler.CreateWorkPaperSignature)
			r.Get("/:id", signatureHandler.GetWorkPaperSignature)
			r.Post("/:id/sign", signatureHandler.SignWorkPaper)
			r.Post("/:id/reject", signatureHandler.RejectWorkPaperSignature)
			r.Post("/:id/reset", signatureHandler.ResetWorkPaperSignature)
		})

		// Work Paper signatures by paper ID
		r.Route("/work-papers/:paperId", func(r fiber.Router) {
			r.Get("/signatures", signatureHandler.GetWorkPaperSignatures)
			r.Get("/pending-signatures", signatureHandler.GetPendingSignaturesByPaperID)
			r.Get("/signature-stats", signatureHandler.GetSignatureStatsByPaperID)
		})

		// User signatures
		r.Route("/users/:userId", func(r fiber.Router) {
			r.Get("/desk/work-papers", signatureHandler.ListWorkPapersWithSignatures)
		})

		// Master LAKIP Item routes (deprecated - for backward compatibility)
		r.Route("/master-lakip-items", func(r fiber.Router) {
			r.Post("/", workPaperItemHandler.CreateMasterLakipItem)
			r.Get("/", workPaperItemHandler.ListMasterLakipItems)
		})

		// Paper Work routes (deprecated - for backward compatibility)
		r.Route("/paper-works", func(r fiber.Router) {
			r.Post("/", workPaperHandler.CreatePaperWork)
		})

		// Paper Work Item routes (deprecated - for backward compatibility)
		r.Route("/paper-work-items/check", func(r fiber.Router) {
			r.Post("/", workPaperHandler.CheckDocument)
		})
	})

	// Legacy routes for backward compatibility
	if businessTripHandler != nil {
		businessTrips := api.Group("/business-trips")
		businessTrips.Use(middleware.AuthMiddleware()) // Apply auth middleware to legacy business trips
		businessTrips.Get("/:id/summary", businessTripHandler.GetBusinessTripSummary)

		// Legacy assignee route
		businessTrips.Post("/:businessTripId/assignees", businessTripHandler.AddAssignee)

		// Legacy transaction routes
		assignees := api.Group("/assignees")
		assignees.Use(middleware.AuthMiddleware()) // Apply auth middleware to legacy assignees
		assignees.Post("/:assigneeId/transactions", businessTripHandler.AddTransaction)
		assignees.Get("/:id/summary", businessTripHandler.GetAssigneeSummary)
	}

	// Cryptography routes
	api.Route("/v1/crypto", func(r fiber.Router) {
		r.Get("/public-key", cryptoHandler.GetPublicKey) // Public key endpoint doesn't need auth

		// Routes that need auth middleware
		r.Use(middleware.AuthMiddleware())
		r.Post("/sign", cryptoHandler.SignDocument)
		r.Post("/verify", cryptoHandler.VerifyDocument)
		r.Post("/verify-offline", cryptoHandler.VerifyDocumentOffline)
		r.Post("/qrcode", cryptoHandler.GenerateQRCode)

		// QR-based signature for work papers
		r.Post("/work-paper-signatures/:signatureId/sign-with-qr", cryptoHandler.SignWorkPaperWithQR)
	})

	// Vaccine routes
	api.Route("/v1/vaccine", func(r fiber.Router) {
		r.Get("/master-vaccines", vaccineHandler.ListMasterVaccines)
		r.Get("/countries", vaccineHandler.ListCountries)
		r.Get("/recommendations/:countryCode", vaccineHandler.GetVaccineRecommendations)
	})

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
		})
	})
}

// Backward compatibility function (deprecated)
func SetupRoutesLegacy(app *fiber.App, transactionHandler *handler.TransactionHandler, meetingHandler *handler.MeetingHandler, businessTripHandler *handler.BusinessTripHandler, assigneeHandler *handler.AssigneeHandler, businessTripTransactionHandler *handler.BusinessTripTransactionHandler, masterLakipItemHandler *deskHandler.WorkPaperItemHandler, paperWorkHandler *deskHandler.WorkPaperHandler) {
	SetupRoutes(app, transactionHandler, meetingHandler, businessTripHandler, assigneeHandler, businessTripTransactionHandler, masterLakipItemHandler, paperWorkHandler, nil, nil, nil)
}
