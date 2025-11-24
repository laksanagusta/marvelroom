package config

import (
	businessTripUC "sandbox/internal/usecase/business_trip"
	meetingUC "sandbox/internal/usecase/meeting"
	transactionUC "sandbox/internal/usecase/transaction"
	"sandbox/internal/infrastructure/drive"
	"sandbox/internal/infrastructure/excel"
	"sandbox/internal/infrastructure/file"
	"sandbox/internal/infrastructure/gemini"
	"sandbox/internal/infrastructure/user"
	postgresInfra "sandbox/internal/infrastructure/database/postgres"
	postgresRepo "sandbox/internal/infrastructure/postgres"
	"sandbox/internal/infrastructure/notification"
	"sandbox/internal/infrastructure/zoom"
	"sandbox/internal/delivery/http/handler"
	"sandbox/internal/domain/repository"
	"sandbox/internal/domain/service"
	"sandbox/pkg/database"

	"github.com/jmoiron/sqlx"
)

// Container holds all application dependencies
type Container struct {
	// Handlers
	TransactionHandler              *handler.TransactionHandler
	MeetingHandler                  *handler.MeetingHandler
	BusinessTripHandler             *handler.BusinessTripHandler
	AssigneeHandler                 *handler.AssigneeHandler
	BusinessTripTransactionHandler  *handler.BusinessTripTransactionHandler

	// Transaction Use Cases
	ExtractTransactionsUseCase *transactionUC.ExtractTransactionsUseCase
	GenerateRecapExcelUseCase  *transactionUC.GenerateRecapExcelUseCase

	// Meeting Use Cases
	CreateMeetingUseCase *meetingUC.CreateMeetingUseCase

	// Business Trip Use Cases
	CreateBusinessTripUseCase                  *businessTripUC.CreateBusinessTripUseCase
	GetBusinessTripUseCase                    *businessTripUC.GetBusinessTripUseCase
	UpdateBusinessTripUseCase                 *businessTripUC.UpdateBusinessTripUseCase
	UpdateBusinessTripWithAssigneesUseCase    *businessTripUC.UpdateBusinessTripWithAssigneesUseCase
	DeleteBusinessTripUseCase                 *businessTripUC.DeleteBusinessTripUseCase
	ListBusinessTripsUseCase                  *businessTripUC.ListBusinessTripsUseCase
	AddAssigneeUseCase                        *businessTripUC.AddAssigneeUseCase
	AddTransactionUseCase                     *businessTripUC.AddTransactionUseCase
	GetBusinessTripSummaryUseCase             *businessTripUC.GetBusinessTripSummaryUseCase
	GetAssigneeSummaryUseCase                 *businessTripUC.GetAssigneeSummaryUseCase

	// Assignee Use Cases
	GetAssigneeUseCase               *businessTripUC.GetAssigneeUseCase
	UpdateAssigneeUseCase            *businessTripUC.UpdateAssigneeUseCase
	DeleteAssigneeUseCase            *businessTripUC.DeleteAssigneeUseCase
	ListAssigneesUseCase             *businessTripUC.ListAssigneesUseCase

	// Transaction Use Cases
	GetTransactionUseCase     *businessTripUC.GetTransactionUseCase
	UpdateTransactionUseCase  *businessTripUC.UpdateTransactionUseCase
	DeleteTransactionUseCase  *businessTripUC.DeleteTransactionUseCase
	ListTransactionsUseCase   *businessTripUC.ListTransactionsUseCase

	// Repositories
	GeminiClient       *gemini.Client
	UserClient         user.ClientInterface
	MeetingRepo        repository.MeetingRepository
	BusinessTripRepo   repository.BusinessTripRepository

	// Database
	DBx   *sqlx.DB

	// Processors
	FileProcessor  *file.Processor
	ExcelGenerator *excel.Generator
}

// NewContainer creates and wires up all dependencies
func NewContainer(cfg *Config) *Container {
	// Initialize database connection using database package
	dbx, err := database.NewConnectionx(cfg.Database.DSN)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// Wrap with database package for consistent interface
	dbWrapper := database.NewDB(dbx)

	// Infrastructure layer
	geminiClient := gemini.NewClient(cfg.Gemini.APIKey)
	userClient := user.NewClientWithAPIKey(cfg.User.BaseURL, cfg.User.APIKey)
	fileProcessor := file.NewProcessor()
	excelGenerator := excel.NewGenerator()

	// Meeting infrastructure
	zoomClient := zoom.NewClient(cfg.Zoom.APIKey, cfg.Zoom.APISecret)
	driveClient := drive.NewClient(cfg.Drive.APIKey)
	notificationClient := notification.NewClient(cfg.Notification.APIKey)
	meetingRepo := postgresInfra.NewRepository(zoomClient, driveClient, notificationClient)

	// Business Trip infrastructure - Now implemented!
	businessTripRepo := postgresRepo.NewBusinessTripRepository(dbWrapper)

	// Business Trip Use Cases - Now enabled!
	createBusinessTripUseCase := businessTripUC.NewCreateBusinessTripUseCase(businessTripRepo, userClient, dbWrapper)
	getBusinessTripUseCase := businessTripUC.NewGetBusinessTripUseCase(businessTripRepo)
	updateBusinessTripUseCase := businessTripUC.NewUpdateBusinessTripUseCase(businessTripRepo)
	updateBusinessTripWithAssigneesUseCase := businessTripUC.NewUpdateBusinessTripWithAssigneesUseCase(businessTripRepo, userClient, dbWrapper)
	deleteBusinessTripUseCase := businessTripUC.NewDeleteBusinessTripUseCase(businessTripRepo)
	listBusinessTripsUseCase := businessTripUC.NewListBusinessTripsUseCase(businessTripRepo)
	addAssigneeUseCase := businessTripUC.NewAddAssigneeUseCase(businessTripRepo, userClient, dbWrapper)
	addTransactionUseCase := businessTripUC.NewAddTransactionUseCase(businessTripRepo)
	getBusinessTripSummaryUseCase := businessTripUC.NewGetBusinessTripSummaryUseCase(businessTripRepo)
	getAssigneeSummaryUseCase := businessTripUC.NewGetAssigneeSummaryUseCase(businessTripRepo)

	// New Assignee Use Cases
	getAssigneeUseCase := businessTripUC.NewGetAssigneeUseCase(businessTripRepo)
	updateAssigneeUseCase := businessTripUC.NewUpdateAssigneeUseCase(businessTripRepo, userClient)
	deleteAssigneeUseCase := businessTripUC.NewDeleteAssigneeUseCase(businessTripRepo)
	listAssigneesUseCase := businessTripUC.NewListAssigneesUseCase(businessTripRepo)

	// New Transaction Use Cases
	getTransactionUseCase := businessTripUC.NewGetTransactionUseCase(businessTripRepo)
	updateTransactionUseCase := businessTripUC.NewUpdateTransactionUseCase(businessTripRepo)
	deleteTransactionUseCase := businessTripUC.NewDeleteTransactionUseCase(businessTripRepo)
	listTransactionsUseCase := businessTripUC.NewListTransactionsUseCase(businessTripRepo)

	// Domain Services
	transactionService := service.NewTransactionService(geminiClient)
	meetingService := service.NewMeetingService(meetingRepo)

	// Transaction Use Cases
	extractTransactionsUseCase := transactionUC.NewExtractTransactionsUseCase(transactionService)
	generateRecapExcelUseCase := transactionUC.NewGenerateRecapExcelUseCase(excelGenerator)

	// Meeting Use Cases
	createMeetingUseCase := meetingUC.NewCreateMeetingUseCase(meetingService)

	// Interface layer
	transactionHandler := handler.NewTransactionHandler(extractTransactionsUseCase, fileProcessor, generateRecapExcelUseCase)
	meetingHandler := handler.NewMeetingHandler(createMeetingUseCase)

	// Business Trip handler - Now enabled!
	businessTripHandler := handler.NewBusinessTripHandler(
		createBusinessTripUseCase,
		getBusinessTripUseCase,
		updateBusinessTripUseCase,
		updateBusinessTripWithAssigneesUseCase,
		deleteBusinessTripUseCase,
		listBusinessTripsUseCase,
		addAssigneeUseCase,
		addTransactionUseCase,
		getBusinessTripSummaryUseCase,
		getAssigneeSummaryUseCase,
	)

	// Assignee handler
	assigneeHandler := handler.NewAssigneeHandler(
		addAssigneeUseCase,
		getAssigneeUseCase,
		updateAssigneeUseCase,
		deleteAssigneeUseCase,
		listAssigneesUseCase,
	)

	// Business Trip Transaction handler
	businessTripTransactionHandler := handler.NewBusinessTripTransactionHandler(
		addTransactionUseCase,
		getTransactionUseCase,
		updateTransactionUseCase,
		deleteTransactionUseCase,
		listTransactionsUseCase,
		getAssigneeUseCase,
	)

	return &Container{
		TransactionHandler:              transactionHandler,
		MeetingHandler:                  meetingHandler,
		BusinessTripHandler:             businessTripHandler,
		AssigneeHandler:                 assigneeHandler,
		BusinessTripTransactionHandler:  businessTripTransactionHandler,
		ExtractTransactionsUseCase:      extractTransactionsUseCase,
		GenerateRecapExcelUseCase:       generateRecapExcelUseCase,
		CreateMeetingUseCase:            createMeetingUseCase,
		CreateBusinessTripUseCase:       createBusinessTripUseCase,
		GetBusinessTripUseCase:          getBusinessTripUseCase,
		UpdateBusinessTripUseCase:       updateBusinessTripUseCase,
		DeleteBusinessTripUseCase:       deleteBusinessTripUseCase,
		ListBusinessTripsUseCase:        listBusinessTripsUseCase,
		AddAssigneeUseCase:              addAssigneeUseCase,
		AddTransactionUseCase:           addTransactionUseCase,
		GetBusinessTripSummaryUseCase:   getBusinessTripSummaryUseCase,
		GetAssigneeSummaryUseCase:       getAssigneeSummaryUseCase,
		GetAssigneeUseCase:              getAssigneeUseCase,
		UpdateAssigneeUseCase:           updateAssigneeUseCase,
		DeleteAssigneeUseCase:           deleteAssigneeUseCase,
		ListAssigneesUseCase:            listAssigneesUseCase,
		GetTransactionUseCase:           getTransactionUseCase,
		UpdateTransactionUseCase:        updateTransactionUseCase,
		DeleteTransactionUseCase:        deleteTransactionUseCase,
		ListTransactionsUseCase:         listTransactionsUseCase,
		GeminiClient:                    geminiClient,
		UserClient:                      userClient,
		MeetingRepo:                     meetingRepo,
		BusinessTripRepo:                businessTripRepo,
		DBx:                             dbx,
		FileProcessor:                   fileProcessor,
		ExcelGenerator:                  excelGenerator,
	}
}
