package repository

import (
	"context"

	"sandbox/internal/domain/entity"
	"sandbox/pkg/pagination"
)

// VaccinesRepository defines the interface for vaccine data operations
type VaccinesRepository interface {
	// Master Vaccine operations
	CreateMasterVaccine(ctx context.Context, vaccine *entity.MasterVaccine) (*entity.MasterVaccine, error)
	GetMasterVaccineByID(ctx context.Context, id string) (*entity.MasterVaccine, error)
	GetMasterVaccineByCode(ctx context.Context, code string) (*entity.MasterVaccine, error)
	UpdateMasterVaccine(ctx context.Context, vaccine *entity.MasterVaccine) (*entity.MasterVaccine, error)
	DeleteMasterVaccine(ctx context.Context, id string) error
	ListMasterVaccines(ctx context.Context, params *pagination.QueryParams) ([]*entity.MasterVaccine, int64, error)
	ListActiveMasterVaccines(ctx context.Context) ([]*entity.MasterVaccine, error)
	ListMasterVaccinesByType(ctx context.Context, vaccineType entity.VaccineType) ([]*entity.MasterVaccine, error)

	// Country operations
	CreateCountry(ctx context.Context, country *entity.Country) (*entity.Country, error)
	GetCountryByID(ctx context.Context, id string) (*entity.Country, error)
	GetCountryByCode(ctx context.Context, code string) (*entity.Country, error)
	UpdateCountry(ctx context.Context, country *entity.Country) (*entity.Country, error)
	DeleteCountry(ctx context.Context, id string) error
	List(ctx context.Context, params *pagination.QueryParams) ([]*entity.Country, int64, error)
	ListActiveCountries(ctx context.Context) ([]*entity.Country, error)
}
