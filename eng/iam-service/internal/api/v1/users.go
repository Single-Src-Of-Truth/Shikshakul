package v1

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Modulix-IT/Shikshakul-WL/iam-service/internal/domain"
	"github.com/Modulix-IT/Shikshakul-WL/iam-service/internal/service"
	"github.com/Modulix-IT/Shikshakul-WL/iam-service/pkg/validator"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// InviteUser handles POST /api/v1/invite
func (h *UserHandler) InviteUser(w http.ResponseWriter, r *http.Request) {
	var req domain.InviteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if errors := validator.ValidateStruct(req); errors != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"status": "error",
			"error":  "Validation Failed",
			"fields": errors,
		})
		return
	}

	// 3. Extract Tenant ID (TODO: In real auth, get this from the JWT/Context)
	// For now, we hardcode it or take it from a header for testing.
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "Missing X-Tenant-ID header", http.StatusUnauthorized)
		return
	}

	token, err := h.userService.InviteUser(r.Context(), tenantID, req)
	if err != nil {
		slog.Error("Failed to invite user", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"status": "success",
		"data": map[string]string{
			"message":      "User invited successfully",
			"invite_token": token, // In prod, do not return this! Send it via email only.
		},
	})
}

// AcceptInvite handles POST /api/v1/accept-invite
func (h *UserHandler) AcceptInvite(w http.ResponseWriter, r *http.Request) {
	var req *domain.AcceptInviteRequest
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

	err := h.userService.AcceptInvite(r.Context(), req.Token, req.Password)
	if err != nil {
		// Security Note: Don't reveal if it was DB error or Token error in too much detail to public
		// But for now, we return the error string.
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Account activated. You can now log in.",
	})
}
