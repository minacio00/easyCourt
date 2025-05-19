package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/minacio00/easyCourt/internal/model"
	"github.com/minacio00/easyCourt/internal/service"
)

type BookingHandler struct {
	service service.BookingService
}

// NewBookingHandler initializes a new BookingHandler
func NewBookingHandler(s service.BookingService) *BookingHandler {
	return &BookingHandler{s}
}

// GetUserBookings godoc
// @Summary      Get user bookings
// @Description  Retrieves bookings for the authenticated user with pagination support
// @Tags         bookings
// @Accept       json
// @Produce      json
// @Param        limit  query int     false "Number of bookings to retrieve (default 10)"
// @Param        offset query int     false "Offset for pagination (default 0)"
// @Success      200    {array}       model.ReadBooking
// @Failure      400    {string}      string "Invalid parameters"
// @Failure      401    {string}      string "Unauthorized"
// @Failure      500    {string}      string "Internal server error"
// @Security     BearerAuth
// @Router       /users/current_user_bookings [get]
func (h *BookingHandler) GetUserBookings(w http.ResponseWriter, r *http.Request) {
	userIDVal := r.Context().Value("user_id")
	if userIDVal == nil {
		http.Error(w, "Unauthorized: user ID not found", http.StatusUnauthorized)
		return
	}
	str, ok := userIDVal.(uint)
	if !ok {
		log.Printf("Invalid user ID type: %T (value: %v)", userIDVal, userIDVal)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	userID := int(str)
	limit := 10
	if l := r.URL.Query().Get("limit"); l != "" {
		parsed, err := strconv.Atoi(l)
		if err != nil || parsed < 0 {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
		limit = parsed
	}

	offset := 0
	if o := r.URL.Query().Get("offset"); o != "" {
		parsed, err := strconv.Atoi(o)
		if err != nil || parsed < 0 {
			http.Error(w, "Invalid offset parameter", http.StatusBadRequest)
			return
		}
		offset = parsed
	}

	bookings, err := h.service.GetUserBookings(userID, limit, offset)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(bookings); err != nil {
		log.Println(err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// CreateBooking godoc
// @Summary      Create a new booking
// @Description  Creates a new booking based on the provided payload
// @Tags         bookings
// @Accept       json
// @Produce      json
// @Param        booking body model.CreateBooking true "Booking Data"
// @Success      201 {object} model.Booking
// @Failure      400 {string} string "Invalid request payload"
// @Failure      500 {string} string "Internal server error"
// @Router       /bookings [post]
func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var createBooking model.CreateBooking
	if err := json.NewDecoder(r.Body).Decode(&createBooking); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateBooking(createBooking.ConvertToBooking()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// GetBookingByID godoc
// @Summary      Get booking by ID
// @Description  Retrieves a specific booking by its ID
// @Tags         bookings
// @Accept       json
// @Produce      json
// @Param        id path int true "Booking ID"
// @Success      200 {object} model.Booking
// @Failure      404 {string} string "Booking not found"
// @Failure      500 {string} string "Internal server error"
// @Router       /bookings/{id} [get]
func (h *BookingHandler) GetBookingByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	booking, err := h.service.GetBookingByID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(booking)
}

// GetAllBookings godoc
// @Summary      Get all bookings
// @Description  Retrieves all bookings with pagination support
// @Tags         bookings
// @Accept       json
// @Produce      json
// @Param        limit  query int false "Number of bookings to retrieve" default(10)
// @Param        offset query int false "Offset of bookings" default(0)
// @Success      200 {array} model.ReadBooking
// @Failure      500 {string} string "Internal server error"
// @Router       /bookings [get]
func (h *BookingHandler) GetAllBookings(w http.ResponseWriter, r *http.Request) {
	// Extract 'limit' and 'offset' from query params
	limitParam := r.URL.Query().Get("limit")
	offsetParam := r.URL.Query().Get("offset")

	// Default values for pagination if not provided
	limit := 0
	offset := 0

	// Parse limit and offset if provided
	if limitParam != "" {
		l, err := strconv.Atoi(limitParam)
		if err == nil {
			limit = l
		}
	}

	if offsetParam != "" {
		o, err := strconv.Atoi(offsetParam)
		if err == nil {
			offset = o
		}
	}

	// Retrieve bookings with pagination
	bookings, err := h.service.GetAllBookings(limit, offset)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return the bookings as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

// UpdateBooking godoc
// @Summary      Update a booking
// @Description  Updates an existing booking based on the provided payload
// @Tags         bookings
// @Accept       json
// @Produce      json
// @Param        id path int true "Booking ID"
// @Param        booking body model.CreateBooking true "Booking Data"
// @Success      200 {object} model.Booking
// @Failure      400 {string} string "Invalid request payload"
// @Failure      404 {string} string "Booking not found"
// @Failure      500 {string} string "Internal server error"
// @Router       /bookings/{id} [put]
func (h *BookingHandler) UpdateBooking(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	var booking model.Booking
	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	booking.ID = id
	if err := h.service.UpdateBooking(&booking); err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(booking)
}

// DeleteBooking godoc
// @Summary      Delete a booking
// @Description  Deletes a booking by its ID
// @Tags         bookings
// @Accept       json
// @Produce      json
// @Param        id path int true "Booking ID"
// @Success      204 {string} string "No Content"
// @Failure      404 {string} string "Booking not found"
// @Failure      500 {string} string "Internal server error"
// @Router       /bookings/{id} [delete]
func (h *BookingHandler) DeleteBooking(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if err := h.service.DeleteBooking(id); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ResetBookings godoc
// @Summary      Delete all the bookings
// @Description  Delete all the bookings
// @Tags         bookings
// @Accept       json
// @Produce      json
// @Success      200 {string} success "ok"
// @Failure      500 {string} string "Internal server error"
// @Router       /bookings/reset [delete]
func (h *BookingHandler) ResetBookings(w http.ResponseWriter, r *http.Request) {
	errMap := make(map[string]string)
	if err := h.service.ResetBookings(); err != nil {
		errMap["error"] = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errMap)
		return
	}
}
