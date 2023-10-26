package inconst

type CtxKey string

const (
	REQID_HEADER = "X-Request-Id"
)

const (
	AUTH_CTX_KEY CtxKey = "auth-ctx"
	MID_CTX_KEY  CtxKey = "mid-ctx"
)

const (
	ROLE_ADMIN    = 1
	ROLE_CUSTOMER = 2
	ROLE_MERCHANT = 3
)
