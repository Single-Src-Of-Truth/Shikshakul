package v1

import (
	"encoding/json"
	"net/http"

	"github.com/Modulix-IT/Shikshakul-WL/iam-service/internal/service"
)

type JwksHandler struct {
	keyManager *service.KeyManager
}

func NewJwksHandler(keyManager *service.KeyManager) *JwksHandler {
	return &JwksHandler{keyManager: keyManager}
}

// GetJWKS handles GET /.well-known/jwks.json
func (h *JwksHandler) GetJWKS(w http.ResponseWriter, r *http.Request) {
	// Allow CORS (Crucial so frontend apps running on other domains can verify tokens)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Cache-Control", "public, max-age=3600") // Cache for 1 hour

	jwks, err := h.keyManager.GetPublicJWKS(r.Context())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jwks)
}
