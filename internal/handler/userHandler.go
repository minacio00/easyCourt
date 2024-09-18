// Package handler manages all user-related operations.
package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/minacio00/easyCourt/internal/model"
	"github.com/minacio00/easyCourt/internal/service"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service}
}

// CreateUser creates a new user
// @Summary Create a new user
// @Description Create a new user with the provided information
// @Tags users
// @Accept  json
// @Produce  json
// @Param   user  body      model.User  true  "User data"
// @Success 201
// @Failure 400  {object}  model.APIError
// @Router /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateUser(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetUserByID retrieves a user by ID
// @Summary Get a user by ID
// @Description Get details of a user by their ID
// @Tags users
// @Produce  json
// @Param   id   path      int  true  "User ID"
// @Success 200  {object}  model.User
// @Failure 404  {object}  model.APIError
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUserByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	resp := user.MapUserToResponse()
	json.NewEncoder(w).Encode(resp)
}

// GetAllUsers retrieves all users
// @Summary Get all users
// @Description Get a list of all users
// @Tags users
// @Produce  json
// @Success 200  {array}  model.User
// @Router /users [get]
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	count := len(users)

	resp := make([]model.UserResponse, count)
	var wg sync.WaitGroup
	wg.Add(count)

	// we don't need to wait for each MapUserToResponse execution to go to the next iteration
	for i, v := range users {
		go func(i int, v *model.User) {
			defer wg.Done()
			resp[i] = *v.MapUserToResponse()
		}(i, &v)
	}
	wg.Wait()

	json.NewEncoder(w).Encode(resp)
}

// UpdateUser updates a user
// @Summary Update a user
// @Description Update the details of a user
// @Tags users
// @Accept  json
// @Produce  json
// @Param   user  body      model.User  true  "Updated user data"
// @Success 204
// @Failure 400  {object}  model.APIError
// @Router /users [put]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateUser(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteUser deletes a user by ID
// @Summary Delete a user by ID
// @Description Delete a user by their ID
// @Tags users
// @Param   id   path      int  true  "User ID"
// @Success 204
// @Failure 404  {object}  model.APIError
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteUser(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
