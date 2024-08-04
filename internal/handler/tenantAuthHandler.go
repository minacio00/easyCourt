// easyCourt/internal/handler/tenant_auth_handler.go
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/minacio00/easyCourt/internal/service"
)

type TenantAuthHandler struct {
	service service.TenantService
}

func NewTenantAuthHandler(service service.TenantService) *TenantAuthHandler {
	return &TenantAuthHandler{service}
}

func (h *TenantAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tenant, err := h.service.Authenticate(credentials.Email, credentials.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(tenant)
}
