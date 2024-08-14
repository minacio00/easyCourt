package handler

import (
	"encoding/json"
	"net/http"

	"github.com/minacio00/easyCourt/internal/service"
)

type UserAuthHandler struct {
	service service.UserService
}

func NewUserAuthHandler(service service.UserService) *UserAuthHandler {
	return &UserAuthHandler{service}
}

func (h *UserAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.Authenticate(credentials.Email, credentials.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(user)
}
