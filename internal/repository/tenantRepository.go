package repository

import (
	"github.com/minacio00/easyCourt/internal/model"
	"gorm.io/gorm"
)

type TenantRepository interface {
	CreateTenant(tenant *model.Tenant) error
	GetTenantByID(id uint) (*model.Tenant, error)
	GetTenantByEmail(email string) (*model.Tenant, error)
	GetAllTenants() ([]model.Tenant, error)
	UpdateTenant(tenant *model.Tenant) error
	DeleteTenant(id uint) error
}

type tenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) TenantRepository {
	return &tenantRepository{db}
}

func (t *tenantRepository) CreateTenant(tenant *model.Tenant) error {
	return t.db.Create(tenant).Error
}
func (r *tenantRepository) GetTenantByID(id uint) (*model.Tenant, error) {
	var tenant model.Tenant
	if err := r.db.First(&tenant, id).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *tenantRepository) GetTenantByEmail(email string) (*model.Tenant, error) {
	var tenant model.Tenant
	if err := r.db.Where("email = ?", email).First(&tenant).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *tenantRepository) GetAllTenants() ([]model.Tenant, error) {
	var tenants []model.Tenant
	if err := r.db.Find(&tenants).Error; err != nil {
		return nil, err
	}
	return tenants, nil
}

func (r *tenantRepository) UpdateTenant(tenant *model.Tenant) error {
	return r.db.Save(tenant).Error
}

func (r *tenantRepository) DeleteTenant(id uint) error {
	return r.db.Delete(&model.Tenant{Model: gorm.Model{ID: id}}).Error
}
