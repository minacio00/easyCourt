package service

import (
	"github.com/minacio00/easyCourt/internal/model"
	"github.com/minacio00/easyCourt/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type TenantService interface {
	CreateTenant(tenant *model.Tenant) error
	GetTenantByID(id uint) (*model.Tenant, error)
	GetTenantByEmail(email string) (*model.Tenant, error)
	GetAllTenants() ([]model.Tenant, error)
	UpdateTenant(tenant *model.Tenant) error
	DeleteTenant(id uint) error
	Authenticate(email, password string) (*model.Tenant, error) // Add this line
}

type tenantService struct {
	repository repository.TenantRepository
}

func NewTenantService(repo repository.TenantRepository) TenantService {
	return &tenantService{repo}
}
func (s *tenantService) CreateTenant(tenant *model.Tenant) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tenant.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	tenant.Password = string(hashedPassword)
	return s.repository.CreateTenant(tenant)
}

func (s *tenantService) GetTenantByID(id uint) (*model.Tenant, error) {
	return s.repository.GetTenantByID(id)
}

func (s *tenantService) GetTenantByEmail(email string) (*model.Tenant, error) {
	return s.repository.GetTenantByEmail(email)
}

func (s *tenantService) GetAllTenants() ([]model.Tenant, error) {
	return s.repository.GetAllTenants()
}

func (s *tenantService) UpdateTenant(tenant *model.Tenant) error {
	return s.repository.UpdateTenant(tenant)
}

func (s *tenantService) DeleteTenant(id uint) error {
	return s.repository.DeleteTenant(id)
}

// need to return jwt
func (s *tenantService) Authenticate(email, password string) (*model.Tenant, error) {
	tenant, err := s.repository.GetTenantByEmail(email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(tenant.Password), []byte(password)); err != nil {
		return nil, err
	}

	return tenant, nil
}
