package indto

type User struct {
	UserID   string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	RoleID   int64  `json:"role_id"`
}
