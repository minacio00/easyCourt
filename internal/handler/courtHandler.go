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
// @Summary Create a new court
// @Description Creates a new court with the provided data
// @Tags court
// @Accept  json
// @Produce  json
// @Param   court  body  model.Court  true  "Court data"
// @Success 201
// @Failure 400  {string}  string  "Invalid request payload"
// @Failure 500  {string}  string  "Failed to create court"
// @Router /courts [post]
func (h *CourtHandler) CreateCourt(w http.ResponseWriter, r *http.Request) {
	var court model.Court
	if err := json.NewDecoder(r.Body).Decode(&court); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := court.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if err := h.service.CreateCourt(&court); err != nil {
		http.Error(w, "Failed to create court", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetAllCourts handles fetching all courts
// @Summary Get all courts
// @Description Retrieves all courts
// @Tags court
// @Produce  json
// @Success 200  {array}  model.Court
// @Failure 500  {string}  string  "Failed to fetch courts"
// @Router /courts [get]
func (h *CourtHandler) GetAllCourts(w http.ResponseWriter, r *http.Request) {
	courts, err := h.service.GetAllCourts()
	if err != nil {
		http.Error(w, "Failed to fetch courts", http.StatusInternalServerError)
		return
	}
	if len(courts) == 0 {
		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode(courts)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courts)
}

// GetCourtByID handles fetching a court by its ID
// @Summary Get court by ID
// @Description Retrieves a court by its ID
// @Tags court
// @Produce  json
// @Param   id  path  int  true  "Court ID"
// @Success 200  {object}  model.Court
// @Failure 400  {string}  string  "Invalid court ID"
// @Failure 404  {string}  string  "Court not found"
// @Router /courts/{id} [get]
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

// @Summary Get courts by location ID
// @Description Retrieves all the courts from a location.
// @Tags court
// @Accept json
// @Produce json
// @Param location_id query int true "ID of the location"
// @Success 200 {array} model.Court "Court information"
// @Failure 400 {object} map[string]string "Error message"
// @Failure 500 {object} map[string]string "Error message"
// @Router /courts/by-location [get]
func (h *CourtHandler) GetCourtByLocation(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("location_id")
	if id == "" {
		http.Error(w, "missing location_id param", http.StatusBadRequest)
		return
	}

	locationID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid location_id param", http.StatusBadRequest)
		return
	}
	courts, err := h.service.GetCourtByLocation(locationID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courts)

}

// UpdateCourt handles updating an existing court
// @Summary Update an existing court
// @Description Updates the data of an existing court
// @Tags court
// @Accept  json
// @Param   court  body  model.Court  true  "Court data"
// @Success 200
// @Failure 400  {string}  string  "Invalid request payload"
// @Failure 500  {string}  string  "Failed to update court"
// @Router /courts [put]
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
// @Summary Delete a court by ID
// @Description Deletes a court based on the given ID
// @Tags court
// @Param   id  path  int  true  "Court ID"
// @Success 204
// @Failure 400  {string}  string  "Invalid court ID"
// @Failure 500  {string}  string  "Failed to delete court"
// @Router /courts/{id} [delete]
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
