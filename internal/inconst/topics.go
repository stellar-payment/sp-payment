package inconst

const (
	// sp-account
	TOPIC_DELETE_USER = "delete-user"

	// sp-payment
	TOPIC_CREATE_MERCHANT = "create-merchant"
	TOPIC_DELETE_MERCHANT = "delete-merchant"
	TOPIC_CREATE_CUSTOMER = "create-customer"
	TOPIC_DELETE_CUSTOMER = "delete-customer"

	// sp-settlement
	TOPIC_CREATE_SETTLEMENT = "create-settlement"
	TOPIC_VOID_SETTLEMENT   = "void-settlement"

	// sp-worker
	TOPIC_CREATE_TRX          = "create-trx"
	TOPIC_CREATE_SCHEDULE_TRX = "create-schedule-trx"
	TOPIC_DELETE_SCHEDULE_TRX = "delete-schedule-trx"
)
