package repository

import (
	"context"

	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/domain"
	"gorm.io/gorm"
)

type TenantRepository interface {
	CreateTenant(ctx context.Context, tenant *domain.Tenant) error
	GetAllTenants(ctx context.Context) ([]domain.Tenant, error)
	GetTenantByID(ctx context.Context, id string) (*domain.Tenant, error)
	UpdateTenant(ctx context.Context, id string, updates map[string]interface{}) error
}

type tenantRepo struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) TenantRepository {
	return &tenantRepo{db: db}
}

func (r *tenantRepo) CreateTenant(ctx context.Context, tenant *domain.Tenant) error {
	return r.db.WithContext(ctx).Create(tenant).Error
}

func (r *tenantRepo) GetAllTenants(ctx context.Context) ([]domain.Tenant, error) {
	var tenants []domain.Tenant
	err := r.db.WithContext(ctx).Order("created_at desc").Find(&tenants).Error
	return tenants, err
}

func (r *tenantRepo) GetTenantByID(ctx context.Context, id string) (*domain.Tenant, error) {
	var tenant domain.Tenant
	if err := r.db.WithContext(ctx).First(&tenant, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *tenantRepo) UpdateTenant(ctx context.Context, id string, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&domain.Tenant{}).Where("id = ?", id).Updates(updates).Error
}
