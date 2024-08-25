package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/minacio00/easyCourt/internal/model"
	"github.com/minacio00/easyCourt/internal/service"
)

type CourtHandler struct {
	service service.CourtService
}

// NewCourtHandler creates a new instance of CourtHandler
func NewCourtHandler(s service.CourtService) *CourtHandler {
	return &CourtHandler{s}
}

// CreateCourt handles the creation of a new court
func (h *CourtHandler) CreateCourt(w http.ResponseWriter, r *http.Request) {
	var court model.Court
	if err := json.NewDecoder(r.Body).Decode(&court); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateCourt(&court); err != nil {
		http.Error(w, "Failed to create court", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetAllCourts handles fetching all courts
func (h *CourtHandler) GetAllCourts(w http.ResponseWriter, r *http.Request) {
	courts, err := h.service.GetAllCourts()
	if err != nil {
		http.Error(w, "Failed to fetch courts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courts)
}

// GetCourtByID handles fetching a court by its ID
func (h *CourtHandler) GetCourtByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid court ID", http.StatusBadRequest)
		return
	}

	court, err := h.service.GetCourtByID(id)
	if err != nil {
		http.Error(w, "Court not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(court)
}

// UpdateCourt handles updating an existing court
func (h *CourtHandler) UpdateCourt(w http.ResponseWriter, r *http.Request) {
	var court model.Court
	if err := json.NewDecoder(r.Body).Decode(&court); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateCourt(&court); err != nil {
		http.Error(w, "Failed to update court", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteCourt handles deleting a court by its ID
func (h *CourtHandler) DeleteCourt(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid court ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteCourt(id); err != nil {
		http.Error(w, "Failed to delete court", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
