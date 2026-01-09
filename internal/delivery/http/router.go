package http

import (
	"sandbox/internal/delivery/http/handler"
	deskHandler "sandbox/internal/delivery/http/handler/desk"
	"sandbox/internal/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App, transactionHandler *handler.TransactionHandler, meetingHandler *handler.MeetingHandler, businessTripHandler *handler.BusinessTripHandler, assigneeHandler *handler.AssigneeHandler, businessTripTransactionHandler *handler.BusinessTripTransactionHandler, workPaperItemHandler *deskHandler.WorkPaperItemHandler, workPaperTopicHandler *deskHandler.WorkPaperTopicHandler, workPaperHandler *deskHandler.WorkPaperHandler, vaccineHandler *handler.VaccineHandler, signatureHandler *handler.WorkPaperSignatureHandler, businessTripDashboardHandler *handler.BusinessTripDashboardHandler, businessTripVerificationHandler *handler.BusinessTripVerificationHandler, grcHandler *handler.GRCHandler, chatbotHandler *handler.ChatbotHandler) {
	api := app.Group("/api")
	api.Post("/upload", middleware.AuthMiddleware(), transactionHandler.UploadAndExtract)
	api.Post("/upload/detailed", middleware.AuthMiddleware(), transactionHandler.UploadAndExtractDetailed)
	api.Post("/report/excel", middleware.AuthMiddleware(), transactionHandler.GenerateRecapExcel)

	api.Post("/meetings", middleware.AuthMiddleware(), meetingHandler.CreateMeeting)

	api.Route("/v1/business-trips", func(r fiber.Router) {
		r.Use(middleware.AuthMiddleware()) // Apply auth middleware to all business trips routes
		r.Get("/dashboard", businessTripDashboardHandler.GetDashboard)
		r.Post("/", businessTripHandler.CreateBusinessTrip)
		r.Get("/", businessTripHandler.ListBusinessTrips)
		r.Get("/verificators", businessTripVerificationHandler.ListVerificators)
		r.Get("/:tripId", businessTripHandler.GetBusinessTrip)
		r.Get("/:tripId/histories", businessTripHandler.GetBusinessTripHistories)
		r.Put("/:tripId", businessTripHandler.UpdateBusinessTrip)
		r.Put("/:tripId/with-assignees", businessTripHandler.UpdateBusinessTripWithAssignees)
		r.Patch("/:tripId/status", businessTripHandler.UpdateBusinessTripStatus)
		r.Delete("/:tripId", businessTripHandler.DeleteBusinessTrip)
		r.Post("/:tripId/verify", businessTripVerificationHandler.VerifyBusinessTrip)

		// Dashboard endpoint
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
			r.Get("/:id", workPaperItemHandler.GetWorkPaperItem)
			r.Put("/:id", workPaperItemHandler.UpdateWorkPaperItem)
			r.Delete("/:id", workPaperItemHandler.DeleteWorkPaperItem)
		})

		// Work Paper Topic routes (new)
		r.Route("/work-paper-topics", func(r fiber.Router) {
			r.Post("/", workPaperTopicHandler.CreateWorkPaperTopic)
			r.Get("/", workPaperTopicHandler.ListWorkPaperTopics)
			r.Get("/active", workPaperTopicHandler.GetActiveWorkPaperTopics)
			r.Get("/:id", workPaperTopicHandler.GetWorkPaperTopic)
			r.Put("/:id", workPaperTopicHandler.UpdateWorkPaperTopic)
			r.Delete("/:id", workPaperTopicHandler.DeleteWorkPaperTopic)
			// Template upload/download/delete routes
			r.Post("/:id/template", workPaperTopicHandler.UploadTemplate)
			r.Get("/:id/template", workPaperTopicHandler.DownloadTemplate)
			r.Delete("/:id/template", workPaperTopicHandler.DeleteTemplate)
		})

		// Work Paper routes (new)
		r.Route("/work-papers", func(r fiber.Router) {
			r.Post("/", workPaperHandler.CreateWorkPaper)
			r.Get("/", workPaperHandler.ListWorkPapers)
			r.Get("/status-transitions", workPaperHandler.GetStatusTransitions)
			r.Get("/:id", workPaperHandler.GetWorkPaperByID)
			r.Put("/:id", workPaperHandler.UpdateWorkPaper)
			r.Put("/:id/status", workPaperHandler.UpdateWorkPaperStatus)
			r.Put("/:id/signers", workPaperHandler.ManageSigners)
			r.Post("/:id/assign-signers", workPaperHandler.AssignSignersBulk)
			r.Get("/:id/docx", workPaperHandler.GenerateDocx)
			r.Get("/:workPaperId/signatures", signatureHandler.GetWorkPaperSignaturesByWorkPaperID)
			// Smart Document Linking - sync folder
			r.Post("/:id/sync-folder", workPaperHandler.SyncFolder)
		})

		// Work Paper Note routes (new)
		r.Post("/work-paper-notes/check", workPaperHandler.CheckWorkPaperNote)
		r.Put("/work-paper-notes/:id", workPaperHandler.UpdateWorkPaperNote)

		// Work Paper Signature routes
		r.Route("/work-paper-signatures", func(r fiber.Router) {
			r.Get("/", signatureHandler.ListWorkPaperSignatures)
			r.Post("/", signatureHandler.CreateWorkPaperSignature)
			r.Get("/:id", signatureHandler.GetWorkPaperSignature)
			r.Post("/:id/sign", signatureHandler.SignWorkPaper)
			r.Post("/:id/reject", signatureHandler.RejectWorkPaperSignature)
			r.Post("/:id/reset", signatureHandler.ResetWorkPaperSignature)
			r.Post("/:id/digital-sign", signatureHandler.CreateDigitalSignature)
			r.Post("/:id/verify", signatureHandler.VerifyDigitalSignature)
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

	// Vaccine routes
	api.Route("/v1/vaccine", func(r fiber.Router) {
		r.Get("/master-vaccines", vaccineHandler.ListMasterVaccines)
		r.Get("/countries", vaccineHandler.ListCountries)
		r.Get("/recommendations/:countryCode", vaccineHandler.GetVaccineRecommendations)
	})

	// GRC Dashboard routes - Live data from Google Spreadsheet
	api.Route("/v1/grc", func(r fiber.Router) {
		r.Get("/overview", grcHandler.GetOverview)
		r.Get("/units", grcHandler.ListUnits)
		r.Get("/units/:id", grcHandler.GetUnitDetail)
		r.Get("/compare", grcHandler.CompareUnits)
		r.Get("/categories", grcHandler.GetCategories)
	})

	// Chatbot routes - RAG-powered chat with file search
	if chatbotHandler != nil {
		api.Route("/v1/chatbot", func(r fiber.Router) {
			r.Use(middleware.AuthMiddleware()) // Apply auth middleware to all chatbot routes

			// Knowledge Base routes
			r.Route("/knowledge-bases", func(r fiber.Router) {
				r.Post("/", chatbotHandler.CreateKnowledgeBase)
				r.Get("/", chatbotHandler.ListKnowledgeBases)
				r.Get("/:id", chatbotHandler.GetKnowledgeBase)
				r.Delete("/:id", chatbotHandler.DeleteKnowledgeBase)
				r.Post("/:id/files", chatbotHandler.UploadFiles)
				r.Post("/:id/sync-status", chatbotHandler.SyncDocumentStatus)
			})

			// Chat Session routes
			r.Route("/sessions", func(r fiber.Router) {
				r.Post("/", chatbotHandler.CreateChatSession)
				r.Get("/", chatbotHandler.ListChatSessions)
				r.Delete("/:id", chatbotHandler.DeleteChatSession)
				r.Post("/:id/messages", chatbotHandler.SendMessage)
				r.Get("/:id/messages", chatbotHandler.GetChatHistory)
			})
		})
	}

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
		})
	})
}

// Backward compatibility function (deprecated)
func SetupRoutesLegacy(app *fiber.App, transactionHandler *handler.TransactionHandler, meetingHandler *handler.MeetingHandler, businessTripHandler *handler.BusinessTripHandler, assigneeHandler *handler.AssigneeHandler, businessTripTransactionHandler *handler.BusinessTripTransactionHandler, masterLakipItemHandler *deskHandler.WorkPaperItemHandler, paperWorkHandler *deskHandler.WorkPaperHandler) {
	SetupRoutes(app, transactionHandler, meetingHandler, businessTripHandler, assigneeHandler, businessTripTransactionHandler, masterLakipItemHandler, nil, paperWorkHandler, nil, nil, nil, nil, nil, nil)
}
