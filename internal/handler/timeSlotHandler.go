package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5" // Assuming you're using chi for routing
	"github.com/minacio00/easyCourt/internal/model"
	"github.com/minacio00/easyCourt/internal/service"
)

type timeSlotHandler struct {
	service service.TimeslotService
}

func NewTimeslotHandler(s service.TimeslotService) *timeSlotHandler {
	return &timeSlotHandler{service: s}
}

// GetTimeslotsByCourt handles fetching all timeslots for a specific court
// @Summary Get timeslots by court
// @Description Retrieves all timeslots for a specific court
// @Tags Timeslots
// @Produce json
// @Param court_id query int true "Court ID"
// @Param day query string false "Weekday filter"
// @Success 200 {array} model.Timeslot "List of timeslots for the court"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "Invalid court ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /timeslots/by-court [get]
func (h *timeSlotHandler) GetTimeslotsByCourt(w http.ResponseWriter, r *http.Request) {
	weekDay := r.URL.Query().Get("day")
	courtIDStr := r.URL.Query().Get("court_id")
	courtID, err := strconv.Atoi(courtIDStr)
	if err != nil {
		http.Error(w, "Invalid court ID", http.StatusBadRequest)
		return
	}

	timeslots, err := h.service.GetTimeslotsByCourt(courtID, weekDay)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(timeslots) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(timeslots)
}

// CreateTimeslot handles the creation of a new timeslot
// @Summary Create a new timeslot
// @Description Creates a new timeslot in the system
// @Tags Timeslots
// @Accept  json
// @Produce  json
// @Param timeslot body model.CreateTimeslot true "Timeslot data"
// @Success 201 {object} map[string]string "Timeslot created successfully"
// @Failure 400 {object} map[string]string "Invalid input data"
// @Router /timeslots [post]
func (h *timeSlotHandler) CreateTimeslot(w http.ResponseWriter, r *http.Request) {
	var timeSlot model.CreateTimeslot
	if err := json.NewDecoder(r.Body).Decode(&timeSlot); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	ts, err := timeSlot.ConvertCreateTimeslotToTimeslot()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	if err := h.service.CreateTimeslot(ts); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Timeslot created successfully"})
}

// GetAllTimeslots handles fetching all timeslots
// @Summary Get all timeslots
// @Description Retrieves all timeslots
// @Tags Timeslots
// @Produce  json
// @Success 200 {array} model.Timeslot "List of timeslots"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /timeslots [get]
func (h *timeSlotHandler) GetAllTimeslots(w http.ResponseWriter, r *http.Request) {
	timeslots, err := h.service.GetAllTimeslots()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if len(timeslots) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(timeslots)
}

// GetTimeslotByID handles fetching a timeslot by its ID
// @Summary Get timeslot by ID
// @Description Retrieves a single timeslot by its ID
// @Tags Timeslots
// @Produce  json
// @Param id path int true "Timeslot ID"
// @Success 200 {object} model.Timeslot "Timeslot data"
// @Failure 400 {object} map[string]string "Invalid ID supplied"
// @Failure 404 {object} map[string]string "Timeslot not found"
// @Router /timeslots/{id} [get]
func (h *timeSlotHandler) GetTimeslotByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid timeslot ID"})
		return
	}

	timeslot, err := h.service.GetTimeslotByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(timeslot)
}

// UpdateTimeslot handles updating an existing timeslot
// @Summary Update a timeslot
// @Description Updates the details of an existing timeslot
// @Tags Timeslots
// @Accept  json
// @Produce  json
// @Param id path int true "Timeslot ID"
// @Param timeslot body model.CreateTimeslot true "Updated timeslot data"
// @Success 200 {object} map[string]string "Timeslot updated successfully"
// @Failure 400 {object} map[string]string "Invalid input or ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /timeslots/{id} [put]
func (h *timeSlotHandler) UpdateTimeslot(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid timeslot ID"})
		return
	}

	var updateTimeslot model.CreateTimeslot
	if err := json.NewDecoder(r.Body).Decode(&updateTimeslot); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	slot, err := updateTimeslot.ConvertCreateTimeslotToTimeslot()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	slot.ID = id // Set the ID from the URL param

	if err := h.service.UpdateTimeslot(slot); err != nil {
		fmt.Println("error during update")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Timeslot updated successfully"})
}

// DeleteTimeslot handles deleting a timeslot by its ID
// @Summary Delete a timeslot
// @Description Deletes a timeslot by its ID
// @Tags Timeslots
// @Param id path int true "Timeslot ID"
// @Success 200 {object} map[string]string "Timeslot deleted successfully"
// @Failure 400 {object} map[string]string "Invalid ID supplied"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /timeslots/{id} [delete]
func (h *timeSlotHandler) DeleteTimeslot(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid timeslot ID"})
		return
	}

	if err := h.service.DeleteTimeslot(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Timeslot deleted successfully"})
}

// GetActiveTimeslots handles fetching all active timeslots
// @Summary Get active timeslots
// @Description Retrieves all active timeslots
// @Tags Timeslots
// @Produce  json
// @Success 200 {array} model.Timeslot "List of active timeslots"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /timeslots/active [get]
func (h *timeSlotHandler) GetActiveTimeslots(w http.ResponseWriter, r *http.Request) {
	timeslots, err := h.service.GetActiveTimeslots()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if len(timeslots) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(timeslots)
}
