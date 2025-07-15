package model 

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponseRepository struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	BackendToken Token   `json:"backend_token"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    uint64 `json:"expires_at"`
}