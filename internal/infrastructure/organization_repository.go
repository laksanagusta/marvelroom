package infrastructure

import (
	"context"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
)

// organizationRepository implements the OrganizationRepository interface
type organizationRepository struct {
	identityService IdentityServiceInterface
}

// NewOrganizationRepository creates a new organization repository
func NewOrganizationRepository(identityService IdentityServiceInterface) repository.OrganizationRepository {
	return &organizationRepository{
		identityService: identityService,
	}
}

func (r *organizationRepository) GetOrganizations(ctx context.Context, page, limit int, sort string) (*entity.OrganizationListResponse, error) {
	return r.identityService.GetOrganizations(ctx, page, limit, sort)
}

func (r *organizationRepository) GetByID(ctx context.Context, id string) (*entity.Organization, error) {
	return r.identityService.GetOrganizationByID(ctx, id)
}