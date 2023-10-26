package indto

type UserResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	RoleID   int64  `json:"role_id"`
}
