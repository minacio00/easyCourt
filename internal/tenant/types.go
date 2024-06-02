package tenant

import "gorm.io/gorm"

type Tenant struct {
	gorm.Model
	Email     string  `json:"email"`
	FirstName string  `json:"firstName"`
	Surname   string  `json:"surname"`
	StripeId  *string `gorm:"->" json:"stripeId"`
	FreeTrial *bool   `json:"freeTrial"`
	Password  *string `json:"password"`
}

type CreateTenantType struct {
	Email     string  `json:"email"`
	FirstName string  `json:"firstName"`
	Surname   string  `json:"surname"`
	FreeTrial *string `json:"freeTrial"`
	Password  *string `json:"password"`
}
