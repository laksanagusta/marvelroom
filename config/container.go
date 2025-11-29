package config

import (
	businessTripUC "sandbox/internal/usecase/business_trip"
	workPaperItemUC "sandbox/internal/usecase/work_paper_item"
	meetingUC "sandbox/internal/usecase/meeting"
	transactionUC "sandbox/internal/usecase/transaction"
	workPaperUC "sandbox/internal/usecase/work_paper"
	vaccineUC "sandbox/internal/usecase/vaccine"
	"sandbox/internal/infrastructure/cdc"
	"sandbox/internal/infrastructure/drive"
	"sandbox/internal/infrastructure/llm"
	"sandbox/internal/infrastructure/excel"
	"sandbox/internal/infrastructure/file"
	"sandbox/internal/infrastructure/gemini"
	"sandbox/internal/infrastructure/cryptography"
	postgresInfra "sandbox/internal/infrastructure/database/postgres"
	postgresRepo "sandbox/internal/infrastructure/postgres"
	"sandbox/internal/infrastructure"
	"sandbox/internal/infrastructure/notification"
	"sandbox/internal/infrastructure/zoom"
	"sandbox/internal/delivery/http/handler"
	deskHandler "sandbox/internal/delivery/http/handler/desk"
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
	WorkPaperItemHandler            *deskHandler.WorkPaperItemHandler
	WorkPaperHandler               *deskHandler.WorkPaperHandler
	WorkPaperSignatureHandler  *handler.WorkPaperSignatureHandler
	VaccineHandler                   *handler.VaccineHandler
	CryptoHandler                    *handler.CryptoHandler

	// Backward compatibility aliases (deprecated)
	MasterLakipItemHandler          *deskHandler.WorkPaperItemHandler
	PaperWorkHandler               *deskHandler.WorkPaperHandler

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

	// Desk Module Use Cases
	CreateWorkPaperItemUseCase *workPaperItemUC.CreateWorkPaperItemUseCase
	ListWorkPaperItemsUseCase  *workPaperItemUC.ListWorkPaperItemsUseCase
	CreateWorkPaperUseCase     *workPaperUC.CreateWorkPaperUseCase
	CheckWorkPaperNoteUseCase   *workPaperUC.CheckWorkPaperNoteUseCase

	// Backward compatibility aliases (deprecated)
	CreateMasterLakipItemUseCase *workPaperItemUC.CreateWorkPaperItemUseCase
	ListMasterLakipItemsUseCase  *workPaperItemUC.ListWorkPaperItemsUseCase
	CreatePaperWorkUseCase       *workPaperUC.CreateWorkPaperUseCase
	CheckDocumentUseCase         *workPaperUC.CheckWorkPaperNoteUseCase

	// Repositories
	GeminiClient            *gemini.Client
	IdentityService         infrastructure.IdentityServiceInterface
	MeetingRepo             repository.MeetingRepository
	BusinessTripRepo        repository.BusinessTripRepository
	WorkPaperItemRepo       repository.WorkPaperItemRepository
	OrganizationRepo        repository.OrganizationRepository
	WorkPaperRepo           repository.WorkPaperRepository
	WorkPaperNoteRepo       repository.WorkPaperNoteRepository
	WorkPaperSignatureRepo repository.WorkPaperSignatureRepository

	// Backward compatibility aliases (deprecated)
	MasterLakipItemRepo     repository.MasterLakipItemRepository
	PaperWorkRepo           repository.PaperWorkRepository
	PaperWorkItemRepo       repository.PaperWorkItemRepository

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
	identityService := infrastructure.NewIdentityServiceWithAPIKey(cfg.User.BaseURL, cfg.User.APIKey)
	fileProcessor := file.NewProcessor()
	excelGenerator := excel.NewGenerator()

	// Meeting infrastructure
	zoomClient := zoom.NewClient(cfg.Zoom.APIKey, cfg.Zoom.APISecret)
	driveClient := drive.NewClient(cfg.Drive.APIKey)
	notificationClient := notification.NewClient(cfg.Notification.APIKey)
	meetingRepo := postgresInfra.NewRepository(zoomClient, driveClient, notificationClient)

	// Business Trip infrastructure - Now implemented!
	businessTripRepo := postgresRepo.NewBusinessTripRepository(dbWrapper)

	// Domain Services - moved up before use cases that use it
	transactionService := service.NewTransactionService(geminiClient)
	meetingService := service.NewMeetingService(meetingRepo)
	userService := service.NewUserService(identityService)
	// vaksinService := service.NewVaksinService(vaksinRepo) // Not used for vaccines endpoint

	// Vaccines infrastructure
	vaccinesRepo := postgresRepo.NewVaccinesRepository(dbWrapper)

	// Business Trip Use Cases - Now enabled!
	createBusinessTripUseCase := businessTripUC.NewCreateBusinessTripUseCase(businessTripRepo, userService, dbWrapper)
	getBusinessTripUseCase := businessTripUC.NewGetBusinessTripUseCase(businessTripRepo)
	updateBusinessTripUseCase := businessTripUC.NewUpdateBusinessTripUseCase(businessTripRepo)
	updateBusinessTripWithAssigneesUseCase := businessTripUC.NewUpdateBusinessTripWithAssigneesUseCase(businessTripRepo, userService, dbWrapper)
	deleteBusinessTripUseCase := businessTripUC.NewDeleteBusinessTripUseCase(businessTripRepo)
	listBusinessTripsUseCase := businessTripUC.NewListBusinessTripsUseCase(businessTripRepo)
	addAssigneeUseCase := businessTripUC.NewAddAssigneeUseCase(businessTripRepo, userService, dbWrapper)
	addTransactionUseCase := businessTripUC.NewAddTransactionUseCase(businessTripRepo)
	getBusinessTripSummaryUseCase := businessTripUC.NewGetBusinessTripSummaryUseCase(businessTripRepo)
	getAssigneeSummaryUseCase := businessTripUC.NewGetAssigneeSummaryUseCase(businessTripRepo)

	// New Assignee Use Cases
	getAssigneeUseCase := businessTripUC.NewGetAssigneeUseCase(businessTripRepo)
	updateAssigneeUseCase := businessTripUC.NewUpdateAssigneeUseCase(businessTripRepo, userService)
	deleteAssigneeUseCase := businessTripUC.NewDeleteAssigneeUseCase(businessTripRepo)
	listAssigneesUseCase := businessTripUC.NewListAssigneesUseCase(businessTripRepo)

	// New Transaction Use Cases
	getTransactionUseCase := businessTripUC.NewGetTransactionUseCase(businessTripRepo)
	updateTransactionUseCase := businessTripUC.NewUpdateTransactionUseCase(businessTripRepo)
	deleteTransactionUseCase := businessTripUC.NewDeleteTransactionUseCase(businessTripRepo)
	listTransactionsUseCase := businessTripUC.NewListTransactionsUseCase(businessTripRepo)
	// CDC Service for vaccine recommendations
	vaccineExtractor := gemini.NewVaccineExtractorAdapter(geminiClient)
	cdcClient := cdc.NewCDCClient(cfg.CDC.BaseURL, cfg.CDC.WebBaseURL, cfg.CDC.APIKey)
	cdcService := service.NewCDCService(vaccinesRepo, cdcClient, vaccineExtractor)

	// Transaction Use Cases
	extractTransactionsUseCase := transactionUC.NewExtractTransactionsUseCase(transactionService)
	generateRecapExcelUseCase := transactionUC.NewGenerateRecapExcelUseCase(excelGenerator)

	// Meeting Use Cases
	createMeetingUseCase := meetingUC.NewCreateMeetingUseCase(meetingService)

	// Vaccines Use Cases
	listMasterVaccinesUseCase := vaccineUC.NewListMasterVaccinesUseCase(vaccinesRepo)
	listCountriesUseCase := vaccineUC.NewListCountriesUseCase(vaccinesRepo)
	getCDCRecommendationsUseCase := vaccineUC.NewGetCDCRecommendationsUseCase(cdcService)

	// Interface layer
	transactionHandler := handler.NewTransactionHandler(extractTransactionsUseCase, fileProcessor, generateRecapExcelUseCase)
	meetingHandler := handler.NewMeetingHandler(createMeetingUseCase)
	vaccineHandler := handler.NewVaccineHandler(listMasterVaccinesUseCase, listCountriesUseCase, getCDCRecommendationsUseCase)

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

	// Desk Module Infrastructure
	workPaperItemRepo := postgresRepo.NewWorkPaperItemRepository(dbWrapper)
	workPaperRepo := postgresRepo.NewWorkPaperRepository(dbWrapper)
	workPaperNoteRepo := postgresRepo.NewWorkPaperNoteRepository(dbWrapper)
	workPaperSignatureRepo := postgresRepo.NewWorkPaperSignatureRepository(dbx)

	// Organization Service - now using unified IdentityService
	organizationRepo := infrastructure.NewOrganizationRepository(identityService)

	// Desk Module Services - Use service account authentication
	gdriveService, err := drive.NewGoogleDriveService("") // Will use default credentials file
	if err != nil {
		panic("Failed to create Google Drive service: " + err.Error())
	}

	llmService, err := llm.NewGeminiService(cfg.Gemini.APIKey)
	if err != nil {
		panic("Failed to create LLM service: " + err.Error())
	}

	deskService := service.NewDeskService(
		workPaperItemRepo,
		organizationRepo,
		workPaperRepo,
		workPaperNoteRepo,
		workPaperSignatureRepo,
		gdriveService,
		llmService,
	)

	// Backward compatibility aliases (deprecated)
	masterLakipItemRepo := workPaperItemRepo
	paperWorkRepo := workPaperRepo
	paperWorkItemRepo := workPaperNoteRepo

	// Desk Module Use Cases
	createWorkPaperItemUseCase := workPaperItemUC.NewCreateWorkPaperItemUseCase(deskService)
	listWorkPaperItemsUseCase := workPaperItemUC.NewListWorkPaperItemsUseCase(workPaperItemRepo)
	createWorkPaperUseCase := workPaperUC.NewCreateWorkPaperUseCase(deskService)
	checkWorkPaperNoteUseCase := workPaperUC.NewCheckWorkPaperNoteUseCase(deskService)
	listWorkPapersUseCase := workPaperUC.NewListWorkPapersUseCase(deskService)
	updateWorkPaperStatusUseCase := workPaperUC.NewUpdateWorkPaperStatusUseCase(deskService)
	updateWorkPaperNoteUseCase := workPaperUC.NewUpdateWorkPaperNoteUseCase(deskService)
	getWorkPaperDetailsUseCase := workPaperUC.NewGetWorkPaperDetailsUseCase(deskService)

	// Backward compatibility aliases
	createMasterLakipItemUseCase := workPaperItemUC.NewCreateMasterLakipItemUseCase(deskService)
	listMasterLakipItemsUseCase := workPaperItemUC.NewListMasterLakipItemsUseCase(workPaperItemRepo)
	createPaperWorkUseCase := workPaperUC.NewCreatePaperWorkUseCase(deskService)
	checkDocumentUseCase := workPaperUC.NewCheckDocumentUseCase(deskService)

	// Desk Module Handlers
	workPaperItemHandler := deskHandler.NewWorkPaperItemHandler(
		createWorkPaperItemUseCase,
		listWorkPaperItemsUseCase,
	)

	workPaperHandler := deskHandler.NewWorkPaperHandler(
		createWorkPaperUseCase,
		checkWorkPaperNoteUseCase,
		listWorkPapersUseCase,
		getWorkPaperDetailsUseCase,
		updateWorkPaperStatusUseCase,
		updateWorkPaperNoteUseCase,
	)

	// Work Paper Signature Handler
	workPaperSignatureHandler := handler.NewWorkPaperSignatureHandler(deskService)

	// Cryptography Service
	cryptoServiceProvider := cryptography.NewServiceProvider()
	cryptoService, err := cryptoServiceProvider.NewCryptoService()
	if err != nil {
		panic("Failed to create crypto service: " + err.Error())
	}

	// Crypto Handler
	cryptoHandler := handler.NewCryptoHandler(cryptoService, deskService)

	// Backward compatibility handler aliases
	masterLakipItemHandler := deskHandler.NewMasterLakipItemHandler(
		createMasterLakipItemUseCase,
		listMasterLakipItemsUseCase,
	)

	paperWorkHandler := deskHandler.NewPaperWorkHandler(
		createPaperWorkUseCase,
		checkDocumentUseCase,
	)

	return &Container{
		TransactionHandler:              transactionHandler,
		MeetingHandler:                  meetingHandler,
		BusinessTripHandler:             businessTripHandler,
		AssigneeHandler:                 assigneeHandler,
		BusinessTripTransactionHandler:  businessTripTransactionHandler,
		WorkPaperItemHandler:            workPaperItemHandler,
		WorkPaperHandler:               workPaperHandler,
		WorkPaperSignatureHandler:  workPaperSignatureHandler,
		VaccineHandler:                   vaccineHandler,
		CryptoHandler:                    cryptoHandler,
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

		// Desk Module Use Cases
		CreateWorkPaperItemUseCase:      createWorkPaperItemUseCase,
		ListWorkPaperItemsUseCase:       listWorkPaperItemsUseCase,
		CreateWorkPaperUseCase:          createWorkPaperUseCase,
		CheckWorkPaperNoteUseCase:        checkWorkPaperNoteUseCase,

		// Backward compatibility aliases (deprecated)
		MasterLakipItemHandler:          masterLakipItemHandler,
		PaperWorkHandler:               paperWorkHandler,
		CreateMasterLakipItemUseCase:    createMasterLakipItemUseCase,
		ListMasterLakipItemsUseCase:     listMasterLakipItemsUseCase,
		CreatePaperWorkUseCase:          createPaperWorkUseCase,
		CheckDocumentUseCase:            checkDocumentUseCase,

		GeminiClient:                    geminiClient,
		IdentityService:                 identityService,
		MeetingRepo:                     meetingRepo,
		BusinessTripRepo:                businessTripRepo,
		WorkPaperItemRepo:               workPaperItemRepo,
		OrganizationRepo:                organizationRepo,
		WorkPaperRepo:                   workPaperRepo,
		WorkPaperNoteRepo:               workPaperNoteRepo,
		WorkPaperSignatureRepo:           workPaperSignatureRepo,

		// Backward compatibility aliases (deprecated)
		MasterLakipItemRepo:             masterLakipItemRepo,
		PaperWorkRepo:                   paperWorkRepo,
		PaperWorkItemRepo:               paperWorkItemRepo,

		DBx:                             dbx,
		FileProcessor:                   fileProcessor,
		ExcelGenerator:                  excelGenerator,
	}
}
