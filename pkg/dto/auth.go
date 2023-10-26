package dto

type AuthLoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthRefreshTokenPayload struct {
	RT string `json:"refresh_token"`
}

type AuthResponse struct {
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	RoleID      int64  `json:"role_id"`
	AccessToken string `json:"access_token"`
}
