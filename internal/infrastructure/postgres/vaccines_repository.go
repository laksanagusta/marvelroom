package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
	"sandbox/pkg/database"
	"sandbox/pkg/pagination"
)

type vaccinesRepository struct {
	db database.Queryer
}

func NewVaccinesRepository(db database.Queryer) repository.VaccinesRepository {
	return &vaccinesRepository{db: db}
}

// Master Vaccine operations

func (r *vaccinesRepository) CreateMasterVaccine(ctx context.Context, vaccine *entity.MasterVaccine) (*entity.MasterVaccine, error) {
	id := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO master_vaccines (
			id, vaccine_code, vaccine_name_id, vaccine_name_en, description_id, description_en,
			vaccine_type, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.ExecContext(ctx, query,
		id, vaccine.VaccineCode, vaccine.VaccineNameID, vaccine.VaccineNameEN,
		vaccine.DescriptionID, vaccine.DescriptionEN, vaccine.VaccineType,
		vaccine.IsActive, now, now,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create master vaccine: %w", err)
	}

	vaccine.ID = id
	vaccine.CreatedAt = now
	vaccine.UpdatedAt = now

	return vaccine, nil
}

func (r *vaccinesRepository) GetMasterVaccineByID(ctx context.Context, id string) (*entity.MasterVaccine, error) {
	query := `
		SELECT id, vaccine_code, vaccine_name_id, vaccine_name_en, description_id, description_en,
			   vaccine_type, is_active, created_at, updated_at
		FROM master_vaccines
		WHERE id = $1
	`

	var vaccine entity.MasterVaccine
	err := r.db.GetContext(ctx, &vaccine, query, id)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("master vaccine not found")
		}
		return nil, fmt.Errorf("failed to get master vaccine: %w", err)
	}

	return &vaccine, nil
}

func (r *vaccinesRepository) GetMasterVaccineByCode(ctx context.Context, code string) (*entity.MasterVaccine, error) {
	query := `
		SELECT id, vaccine_code, vaccine_name_id, vaccine_name_en, description_id, description_en,
			   vaccine_type, is_active, created_at, updated_at
		FROM master_vaccines
		WHERE vaccine_code = $1
	`

	var vaccine entity.MasterVaccine
	err := r.db.GetContext(ctx, &vaccine, query, code)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("master vaccine not found")
		}
		return nil, fmt.Errorf("failed to get master vaccine: %w", err)
	}

	return &vaccine, nil
}

func (r *vaccinesRepository) UpdateMasterVaccine(ctx context.Context, vaccine *entity.MasterVaccine) (*entity.MasterVaccine, error) {
	now := time.Now()
	query := `
		UPDATE master_vaccines
		SET vaccine_name_id = $2, vaccine_name_en = $3, description_id = $4, description_en = $5,
			vaccine_type = $6, is_active = $7, updated_at = $8
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query,
		vaccine.ID, vaccine.VaccineNameID, vaccine.VaccineNameEN,
		vaccine.DescriptionID, vaccine.DescriptionEN, vaccine.VaccineType,
		vaccine.IsActive, now,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update master vaccine: %w", err)
	}

	vaccine.UpdatedAt = now
	return vaccine, nil
}

func (r *vaccinesRepository) DeleteMasterVaccine(ctx context.Context, id string) error {
	query := `DELETE FROM master_vaccines WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete master vaccine: %w", err)
	}

	return nil
}

func (r *vaccinesRepository) ListMasterVaccines(ctx context.Context, params *pagination.QueryParams) ([]*entity.MasterVaccine, int64, error) {
	// Build count query
	countBuilder := pagination.NewQueryBuilder("SELECT COUNT(*) FROM master_vaccines")
	for _, filter := range params.Filters {
		if err := countBuilder.AddFilter(filter); err != nil {
			return nil, 0, err
		}
	}

	countQuery, countArgs := countBuilder.Build()
	var totalCount int64
	err := r.db.GetContext(ctx, &totalCount, countQuery, countArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count master vaccines: %w", err)
	}

	// Build main query with pagination
	queryBuilder := pagination.NewQueryBuilder(`
		SELECT
			id, vaccine_code, vaccine_name_id, vaccine_name_en, description_id, description_en,
			vaccine_type, is_active, created_at, updated_at
		FROM master_vaccines`)

	// Add filters
	for _, filter := range params.Filters {
		if err := queryBuilder.AddFilter(filter); err != nil {
			return nil, 0, err
		}
	}

	// Add sorts
	for _, sort := range params.Sorts {
		if err := queryBuilder.AddSort(sort); err != nil {
			return nil, 0, err
		}
	}

	// Add default sort if no sorts provided
	if len(params.Sorts) == 0 {
		queryBuilder.AddSort(pagination.Sort{
			Field: "vaccine_code",
			Order: "asc",
		})
	}

	query, args := queryBuilder.Build()

	// Add pagination
	offset := (params.Pagination.Page - 1) * params.Pagination.Limit
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", params.Pagination.Limit, offset)

	var vaccines []*entity.MasterVaccine
	err = r.db.SelectContext(ctx, &vaccines, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list master vaccines: %w", err)
	}

	return vaccines, totalCount, nil
}

func (r *vaccinesRepository) ListActiveMasterVaccines(ctx context.Context) ([]*entity.MasterVaccine, error) {
	query := `
		SELECT id, vaccine_code, vaccine_name_id, vaccine_name_en, description_id, description_en,
			   vaccine_type, is_active, created_at, updated_at
		FROM master_vaccines
		WHERE is_active = true
		ORDER BY vaccine_code ASC
	`

	var vaccines []*entity.MasterVaccine
	err := r.db.SelectContext(ctx, &vaccines, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list active master vaccines: %w", err)
	}

	return vaccines, nil
}

func (r *vaccinesRepository) ListMasterVaccinesByType(ctx context.Context, vaccineType entity.VaccineType) ([]*entity.MasterVaccine, error) {
	query := `
		SELECT id, vaccine_code, vaccine_name_id, vaccine_name_en, description_id, description_en,
			   vaccine_type, is_active, created_at, updated_at
		FROM master_vaccines
		WHERE vaccine_type = $1 AND is_active = true
		ORDER BY vaccine_code ASC
	`

	var vaccines []*entity.MasterVaccine
	err := r.db.SelectContext(ctx, &vaccines, query, vaccineType)
	if err != nil {
		return nil, fmt.Errorf("failed to list master vaccines by type: %w", err)
	}

	return vaccines, nil
}

// Country operations

func (r *vaccinesRepository) CreateCountry(ctx context.Context, country *entity.Country) (*entity.Country, error) {
	id := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO countries (
			id, country_code, country_name_id, country_name_en, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		id, country.CountryCode, country.CountryNameID, country.CountryNameEN,
		country.IsActive, now, now,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create country: %w", err)
	}

	country.ID = id
	country.CreatedAt = now
	country.UpdatedAt = now

	return country, nil
}

func (r *vaccinesRepository) GetCountryByID(ctx context.Context, id string) (*entity.Country, error) {
	query := `
		SELECT id, country_code, country_name_id, country_name_en, is_active, created_at, updated_at
		FROM countries
		WHERE id = $1
	`

	var country entity.Country
	err := r.db.GetContext(ctx, &country, query, id)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("country not found")
		}
		return nil, fmt.Errorf("failed to get country: %w", err)
	}

	return &country, nil
}

func (r *vaccinesRepository) GetCountryByCode(ctx context.Context, code string) (*entity.Country, error) {
	query := `
		SELECT id, country_code, country_name_id, country_name_en, is_active, created_at, updated_at
		FROM countries
		WHERE country_code = $1
	`

	var country entity.Country
	err := r.db.GetContext(ctx, &country, query, code)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("country not found")
		}
		return nil, fmt.Errorf("failed to get country: %w", err)
	}

	return &country, nil
}

func (r *vaccinesRepository) UpdateCountry(ctx context.Context, country *entity.Country) (*entity.Country, error) {
	now := time.Now()
	query := `
		UPDATE countries
		SET country_name_id = $2, country_name_en = $3, is_active = $4, updated_at = $5
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query,
		country.ID, country.CountryNameID, country.CountryNameEN,
		country.IsActive, now,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update country: %w", err)
	}

	country.UpdatedAt = now
	return country, nil
}

func (r *vaccinesRepository) DeleteCountry(ctx context.Context, id string) error {
	query := `DELETE FROM countries WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete country: %w", err)
	}

	return nil
}

func (r *vaccinesRepository) List(ctx context.Context, params *pagination.QueryParams) ([]*entity.Country, int64, error) {
	// Build count query
	countBuilder := pagination.NewQueryBuilder("SELECT COUNT(*) FROM countries")
	for _, filter := range params.Filters {
		if err := countBuilder.AddFilter(filter); err != nil {
			return nil, 0, err
		}
	}

	countQuery, countArgs := countBuilder.Build()
	var totalCount int64
	err := r.db.GetContext(ctx, &totalCount, countQuery, countArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count countries: %w", err)
	}

	// Build main query with pagination
	queryBuilder := pagination.NewQueryBuilder(`
		SELECT
			id, country_code, country_name_id, country_name_en, is_active, created_at, updated_at
		FROM countries`)

	// Add filters
	for _, filter := range params.Filters {
		if err := queryBuilder.AddFilter(filter); err != nil {
			return nil, 0, err
		}
	}

	// Add sorts
	for _, sort := range params.Sorts {
		if err := queryBuilder.AddSort(sort); err != nil {
			return nil, 0, err
		}
	}

	// Add default sort if no sorts provided
	if len(params.Sorts) == 0 {
		queryBuilder.AddSort(pagination.Sort{
			Field: "country_name_en",
			Order: "asc",
		})
	}

	query, args := queryBuilder.Build()

	// Add pagination
	offset := (params.Pagination.Page - 1) * params.Pagination.Limit
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", params.Pagination.Limit, offset)

	var countries []*entity.Country
	err = r.db.SelectContext(ctx, &countries, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list countries: %w", err)
	}

	return countries, totalCount, nil
}

func (r *vaccinesRepository) ListActiveCountries(ctx context.Context) ([]*entity.Country, error) {
	query := `
		SELECT id, country_code, country_name_id, country_name_en, is_active, created_at, updated_at
		FROM countries
		WHERE is_active = true
		ORDER BY country_name_en ASC
	`

	var countries []*entity.Country
	err := r.db.SelectContext(ctx, &countries, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list active countries: %w", err)
	}

	return countries, nil
}


