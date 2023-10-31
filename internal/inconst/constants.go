package inconst

type CtxKey string

const (
	REQID_HEADER = "X-Request-Id"
)

const (
	AUTH_CTX_KEY  CtxKey = "auth-ctx"
	TOKEN_CTX_KEY CtxKey = "token-ctx"
	MID_CTX_KEY   CtxKey = "mid-ctx"
)

const (
	ROLE_ADMIN    = 1
	ROLE_CUSTOMER = 2
	ROLE_MERCHANT = 3
)

const (
	ACCOUNT_TYPE_CUST     = 1
	ACCOUNT_TYPE_MERCHANT = 2
)
