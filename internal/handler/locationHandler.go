package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/minacio00/easyCourt/internal/model"
	"github.com/minacio00/easyCourt/internal/service"
)

type LocationHandler struct {
	service service.LocationService
}

func NewLocationHandler(s service.LocationService) *LocationHandler {
	return &LocationHandler{s}
}

// CreateLocation creates a new location
// @Summary Create a new location
// @Description Create a new location with the provided information
// @Tags location
// @Accept  json
// @Produce  json
// @Param   location  body      model.CreateLocation  true  "Location data"
// @Success 201  {object}  model.Location
// @Failure 400  {object}  model.APIError
// @Router /location [post]
func (h *LocationHandler) CreateLocation(w http.ResponseWriter, r *http.Request) {
	var location model.Location
	if err := json.NewDecoder(r.Body).Decode(&location); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	if err := h.service.CreateLocation(&location); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/location/%d", location.ID))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(location)
}

// GetAllLocations retrieves all locations
// @Summary Get all locations
// @Description Get a list of all locations
// @Tags location
// @Produce  json
// @Success 200  {array}  model.Location
// @Success 204 {string} string "No Content"
// @Router /location [get]
func (h *LocationHandler) GetAllLocations(w http.ResponseWriter, r *http.Request) {
	locations, err := h.service.GetAllLocations()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	if len(locations) == 0 {
		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode(locations)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&locations)
}

// UpdateLocation updates an existing location
// @Summary Update a location
// @Description Update the details of a location
// @Tags location
// @Accept  json
// @Produce  json
// @Param   location  body      model.Location  true  "Updated location data"
// @Success 204
// @Failure 400  {object}  model.APIError
// @Router /location [put]
func (h *LocationHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	var location model.Location
	if err := json.NewDecoder(r.Body).Decode(&location); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	if err := location.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	err := h.service.UpdateLocation(&location)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeleteLocation deletes a location by ID
// @Summary Delete a location by ID
// @Description Delete a location by its ID
// @Tags location
// @Param   id   path      int  true  "Location ID"
// @Success 204
// @Failure 400  {object}  model.APIError
// @Router /location/{id} [delete]
func (h *LocationHandler) DeleteLocation(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteLocation(uint(id)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
