package v1

import (
	"context"
	"encoding/json"
	"net"
	"net/http"

	"github.com/Modulix-IT/Shikshakul-WL/iam-service/internal/domain"
	"github.com/Modulix-IT/Shikshakul-WL/iam-service/internal/service"
	"github.com/Modulix-IT/Shikshakul-WL/iam-service/pkg/validator"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login handles POST /api/v1/oauth/token
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if errors := validator.ValidateStruct(req); errors != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"error":  "Validation Failed",
			"fields": errors,
		})
		return
	}

	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "Missing X-Tenant-ID", http.StatusBadRequest)
		return
	}

	ctx := context.WithValue(r.Context(), "tenant_id", tenantID)
	userAgent := r.Header.Get("User-Agent")
	ip := r.RemoteAddr
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		ip = host
	} else {
		ip = r.RemoteAddr
	}

	resp, err := h.authService.Login(ctx, req, userAgent, ip)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Refresh handles POST /api/v1/oauth/refresh
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req domain.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if errors := validator.ValidateStruct(req); errors != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"error":  "Validation Failed",
			"fields": errors,
		})
		return
	}

	resp, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Revoke handles POST /api/v1/oauth/revoke (Logout)
func (h *AuthHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	var req domain.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	_ = h.authService.Logout(r.Context(), req.RefreshToken)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}
