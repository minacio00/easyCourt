package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-chi/chi/v5"
	"github.com/minacio00/easyCourt/internal/model"
	"github.com/minacio00/easyCourt/internal/service"
	"github.com/spf13/viper"
)

type LocationHandler struct {
	service  service.LocationService
	s3Client *s3.S3
	s3Bucket string
}

func NewLocationHandler(s service.LocationService) *LocationHandler {
	secret := viper.GetString("AWS_SECRET_ACCESS_KEY")
	id := viper.GetString("AWS_ACCESS_KEY_ID")
	var creds *credentials.Credentials

	if id != "" && secret != "" {
		creds = credentials.NewStaticCredentials(id, secret, "")
	}
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: creds,
	})
	if err != nil {
		panic(err)
	}

	return &LocationHandler{
		service:  s,
		s3Client: s3.New(sess),
		s3Bucket: "easycourt-locations",
	}
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

// UploadLocationImage uploads an image for a specific location
// @Summary Upload a location image
// @Description Upload an image file for a specific location and store it in S3
// @Tags location
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Location ID"
// @Param image formData file true "Image file to upload"
// @Success 200 {object} map[string]string "Returns a success message and the URL of the uploaded image"
// @Failure 400 {object} model.APIError "Bad request"
// @Failure 404 {object} model.APIError "Location not found"
// @Failure 500 {object} model.APIError "Internal server error"
// @Router /location/{id}/image [post]
func (h *LocationHandler) UploadLocationImage(w http.ResponseWriter, r *http.Request) {
	// Parse the location ID from the URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid location ID", http.StatusBadRequest)
		return
	}

	// Check if the location exists
	location, err := h.service.GetLocationById(uint(id))
	if err != nil {
		http.Error(w, "Location not found", http.StatusNotFound)
		return
	}

	// Parse the multipart form
	err = r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get the file from the form data
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Unable to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Generate a unique filename

	filename := fmt.Sprintf("location_%d%s", id, filepath.Ext(header.Filename))

	// Upload the file to S3
	_, err = h.s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(h.s3Bucket),
		Key:    aws.String(filename),
		Body:   file,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Update the location with the new image URL
	imageURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", h.s3Bucket, filename)
	location.ImageUrl = imageURL
	err = h.service.UpdateLocation(location)
	if err != nil {
		http.Error(w, "Unable to update location with new image URL", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Image uploaded successfully", "url": imageURL})
}
