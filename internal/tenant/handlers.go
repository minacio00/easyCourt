package tenant

import (
	"encoding/json"
	"net/http"

	"github.com/minacio00/easyCourt/internal"
)

func CreateTenantHandler(w http.ResponseWriter, r *http.Request) {
	var tenant *CreateTenantType
	tnService := &TenantService{}

	err := json.NewDecoder(r.Body).Decode(&tenant)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Check for missing fields

	if err := internal.CheckMissingFields(tenant); err != nil {
		internal.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	newTn := tnService.CreateTenant(tenant)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTn)
}
