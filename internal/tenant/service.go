package tenant

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

type TenantService struct {
}

func (tn *TenantService) CreateTenant(data *CreateTenantType) *Tenant {
	trial := new(bool)
	*trial = true

	hash, err := bcrypt.GenerateFromPassword([]byte(*data.Password), bcrypt.MaxCost)
	hashedPassword := string(hash)
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
	tenant := &Tenant{
		Email:     data.Email,
		FirstName: data.FirstName,
		Surname:   data.Surname,
		FreeTrial: trial,
		Password:  &hashedPassword,
	}

	return tenant
}
