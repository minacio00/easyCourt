// easyCourt/internal/handler/tenant_handler.go
package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/minacio00/easyCourt/internal/model"
	"github.com/minacio00/easyCourt/internal/service"
)

type TenantHandler struct {
	service service.TenantService
}

func NewTenantHandler(service service.TenantService) *TenantHandler {
	return &TenantHandler{service}
}

func (h *TenantHandler) CreateTenant(w http.ResponseWriter, r *http.Request) {
	var tenant model.Tenant
	if err := json.NewDecoder(r.Body).Decode(&tenant); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateTenant(&tenant); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tenant)
}

func (h *TenantHandler) GetTenantByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	tenant, err := h.service.GetTenantByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(tenant)
}

func (h *TenantHandler) GetAllTenants(w http.ResponseWriter, r *http.Request) {
	tenants, err := h.service.GetAllTenants()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tenants)
}

func (h *TenantHandler) UpdateTenant(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var tenant model.Tenant
	if err := json.NewDecoder(r.Body).Decode(&tenant); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tenant.ID = uint(id)
	if err := h.service.UpdateTenant(&tenant); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tenant)
}

func (h *TenantHandler) DeleteTenant(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteTenant(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
