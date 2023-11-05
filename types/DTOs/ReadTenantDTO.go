package dtos

type ReadTenantDTO struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	TrialPeriod bool   `json:"periodo_teste"`
}
