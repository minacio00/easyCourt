package tenant

import "gorm.io/gorm"

type Tenant struct {
	gorm.Model
	Email     string
	FirstName string
	Surname   string
	StripeId  *string
	FreeTrial *string
	Password  *string
}
