package domain

// LoginRequest matches standard OAuth 2.0 parameters
type LoginRequest struct {
	GrantType string `json:"grant_type" validate:"required,eq=password"`
	Username  string `json:"username" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	Scope     string `json:"scope"`
}

// TokenResponse is the standard OAuth 2.0 JSON
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
