package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/your-org/driftlock/api-server/internal/models"
	"github.com/your-org/driftlock/api-server/internal/services"
)

// TenantHandler handles tenant-related HTTP requests
type TenantHandler struct {
	tenantService *services.TenantService
}

// NewTenantHandler creates a new tenant handler
func NewTenantHandler(tenantService *services.TenantService) *TenantHandler {
	return &TenantHandler{
		tenantService: tenantService,
	}
}

// CreateTenant handles POST /tenants
func (h *TenantHandler) CreateTenant(w http.ResponseWriter, r *http.Request) {
	var tenant models.Tenant
	if err := json.NewDecoder(r.Body).Decode(&tenant); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdTenant, err := h.tenantService.CreateTenant(r.Context(), &tenant)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTenant)
}

// GetTenant handles GET /tenants/{id}
func (h *TenantHandler) GetTenant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["id"]

	tenant, err := h.tenantService.GetTenant(r.Context(), tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenant)
}

// UpdateTenant handles PUT /tenants/{id}
func (h *TenantHandler) UpdateTenant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["id"]

	var tenant models.Tenant
	if err := json.NewDecoder(r.Body).Decode(&tenant); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tenant.ID = tenantID
	updatedTenant, err := h.tenantService.UpdateTenant(r.Context(), &tenant)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTenant)
}

// DeleteTenant handles DELETE /tenants/{id}
func (h *TenantHandler) DeleteTenant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["id"]

	err := h.tenantService.DeleteTenant(r.Context(), tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListTenants handles GET /tenants
func (h *TenantHandler) ListTenants(w http.ResponseWriter, r *http.Request) {
	tenants, err := h.tenantService.ListTenants(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenants)
}

// GetTenantUsage handles GET /tenants/{id}/usage
func (h *TenantHandler) GetTenantUsage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["id"]

	tenant, err := h.tenantService.GetTenant(r.Context(), tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenant.Usage)
}

// UpdateTenantUsage handles POST /tenants/{id}/usage
func (h *TenantHandler) UpdateTenantUsage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["id"]

	var request struct {
		Type   string `json:"type"`
		Amount int    `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.tenantService.UpdateTenantUsage(r.Context(), tenantID, request.Type, request.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CheckTenantQuota handles GET /tenants/{id}/quota/{resource}
func (h *TenantHandler) CheckTenantQuota(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["id"]
	resourceType := vars["resource"]

	exceeded, err := h.tenantService.CheckTenantQuota(r.Context(), tenantID, resourceType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"exceeded": exceeded,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListTenantsWithPagination handles GET /tenants with pagination
func (h *TenantHandler) ListTenantsWithPagination(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// Get all tenants (simplified - in production, you'd implement pagination in the service)
	tenants, err := h.tenantService.ListTenants(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Calculate pagination
	total := len(tenants)
	totalPages := (total + pageSize - 1) / pageSize
	start := (page - 1) * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}

	// Slice the results
	if start >= total {
		tenants = []*models.Tenant{}
	} else {
		if end > total {
			end = total
		}
		tenants = tenants[start:end]
	}

	response := map[string]interface{}{
		"tenants":     tenants,
		"page":        page,
		"page_size":   pageSize,
		"total":       total,
		"total_pages": totalPages,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
