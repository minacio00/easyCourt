package handler

import (
	"encoding/json"
	"net/http"

	_ "github.com/minacio00/easyCourt/internal/model"
	"github.com/minacio00/easyCourt/internal/service"
)

type UserAuthHandler struct {
	service service.UserService
}

func NewUserAuthHandler(service service.UserService) *UserAuthHandler {
	return &UserAuthHandler{service}
}

type Credentials struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type Refresh struct {
	RefreshToken string `json:"refresh_token"`
}

// @Summary User login
// @Description Authenticate a user and return access and refresh tokens
// @Tags authentication
// @Accept json
// @Produce json
// @Param credentials body Credentials true "Login Credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} model.APIError
// @Failure 401 {object} model.APIError
// @Router /login [post]
func (h *UserAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, token, refresh, err := h.service.Authenticate(credentials.Phone, credentials.Password)
	if err != nil {
		http.Error(w, "Invalid phone or password", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  token,
		"refresh_token": refresh,
	})
}

// Refresh godoc
// @Summary Refresh access token
// @Description Use a refresh token to obtain a new access token and refresh token pair
// @Tags authentication
// @Accept json
// @Produce json
// @Param refresh body Refresh true "Refresh Token"
// @Success 200 {object} map[string]string
// @Failure 400 {object} model.APIError
// @Failure 401 {object} model.APIError
// @Router /refresh [post]
func (h *UserAuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var refreshReq map[string]string
	if err := json.NewDecoder(r.Body).Decode(&refreshReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	refreshToken, ok := refreshReq["refresh_token"]
	if !ok {
		http.Error(w, "Refresh token is required", http.StatusBadRequest)
		return
	}

	newAccessToken, newRefreshToken, err := h.service.RefreshToken(refreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}
