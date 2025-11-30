package business_trip

import (
	"context"

	"sandbox/internal/domain/repository"
	"sandbox/pkg/pagination"
)

type ListVerificatorsUseCase struct {
	businessTripRepo repository.BusinessTripRepository
}

func NewListVerificatorsUseCase(businessTripRepo repository.BusinessTripRepository) *ListVerificatorsUseCase {
	return &ListVerificatorsUseCase{
		businessTripRepo: businessTripRepo,
	}
}

// ListVerificatorsResponse represents the response for listing verificators
type ListVerificatorsResponse struct {
	ID                string         `json:"id"`
	BusinessTripID    string         `json:"business_trip_id"`
	UserID            string         `json:"user_id"`
	UserName          string         `json:"user_name"`
	EmployeeNumber    string         `json:"employee_number"`
	Position          string         `json:"position"`
	Status            string         `json:"status"`
	VerificationNotes string         `json:"verification_notes"`
	VerifiedAt        string         `json:"verified_at,omitempty"`
	BusinessTrip      *BusinessTrip  `json:"business_trip"`
}

// BusinessTrip represents business trip data for list verificators response
type BusinessTrip struct {
	ID                    string `json:"id"`
	BusinessTripNumber    string `json:"business_trip_number"`
	ActivityPurpose       string `json:"activity_purpose"`
	DestinationCity       string `json:"destination_city"`
	StartDate             string `json:"start_date"`
	EndDate               string `json:"end_date"`
	SPDDate               string `json:"spd_date"`
	DepartureDate         string `json:"departure_date"`
	ReturnDate            string `json:"return_date"`
	Status                string `json:"status"`
	DocumentLink          string `json:"document_link,omitempty"`
}

func (uc *ListVerificatorsUseCase) Execute(ctx context.Context, params *pagination.QueryParams) ([]*ListVerificatorsResponse, *pagination.PagedResponse, error) {
	// Get verificators with pagination directly from repository
	verificators, totalCount, err := uc.businessTripRepo.ListVerificators(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	// Convert entities to response DTOs
	var responses []*ListVerificatorsResponse
	for _, v := range verificators {
		response := &ListVerificatorsResponse{
			ID:                v.ID,
			BusinessTripID:    v.BusinessTripID,
			UserID:            v.UserID,
			UserName:          v.UserName,
			EmployeeNumber:    v.EmployeeNumber,
			Position:          v.Position,
			Status:            string(v.Status),
			VerificationNotes: v.VerificationNotes,
			BusinessTrip: &BusinessTrip{
				ID:                 v.BusinessTripID,
				ActivityPurpose:    v.BusinessTripActivityPurpose,
				DestinationCity:     v.BusinessTripDestinationCity,
				StartDate:           v.BusinessTripStartDate.Format("2006-01-02"),
				EndDate:             v.BusinessTripEndDate.Format("2006-01-02"),
				SPDDate:             v.BusinessTripSPDDate.Format("2006-01-02"),
				DepartureDate:       v.BusinessTripDepartureDate.Format("2006-01-02"),
				ReturnDate:          v.BusinessTripReturnDate.Format("2006-01-02"),
				Status:              string(v.BusinessTripStatus),
			},
		}

		// Handle BusinessTripNumber which is sql.NullString
		if v.BusinessTripNumber.Valid {
			response.BusinessTrip.BusinessTripNumber = v.BusinessTripNumber.String
		}

		// Handle DocumentLink which is sql.NullString
		if v.BusinessTripDocumentLink.Valid {
			response.BusinessTrip.DocumentLink = v.BusinessTripDocumentLink.String
		}

		if v.VerifiedAt != nil {
			response.VerifiedAt = v.VerifiedAt.Format("2006-01-02T15:04:05Z07:00")
		}

		responses = append(responses, response)
	}

	// Calculate pagination info
	totalPages := int(totalCount) / params.Pagination.Limit
	if int(totalCount)%params.Pagination.Limit > 0 {
		totalPages++
	}

	return responses, &pagination.PagedResponse{
		Page:       params.Pagination.Page,
		Limit:      params.Pagination.Limit,
		TotalItems: totalCount,
		TotalPages: totalPages,
	}, nil
}
