package inconst

type CtxKey string

const (
	REQID_HEADER     = "X-Request-Id"
	CORRREQID_HEADER = "X-Correlation-Id"
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

const (
	TRX_TYPE_P2P         = 1
	TRX_TYPE_P2B         = 2
	TRX_TYPE_BENEFICIARY = 3
	TRX_TYPE_SYSTEM      = 9
)

const (
	TRX_STATUS_PENDING   = 0
	TRX_STATUS_SUCCESS   = 1
	TRX_STATUS_CANCELLED = 2
	TRX_STATUS_VOID      = 9
)

const (
	BNF_STATUS_PENDING = 0
	BNF_STATUS_CONFIRM = 1
)
