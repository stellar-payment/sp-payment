package router

const (
	basePath = "/payment/api/v1"
	PingPath = basePath + "/ping"

	// ----- Customers
	customerBasepath = basePath + "/customers"
	customerIDPath   = customerBasepath + "/:customerID"

	// ----- Merchants
	merchantBasepath = basePath + "/merchants"
	merchantIDPath   = merchantBasepath + "/:merchantID"

	// ----- Accounts
	accountBasepath = basePath + "/accounts"
	accountIDPath   = accountBasepath + "/:accountID"

	// ----- Transactions
	trxBasepath = basePath + "/transactions"
	trxIDPath   = trxBasepath + "/:trxID"
	trxP2PPath  = trxBasepath + "/p2p"
	trxP2BPath  = trxBasepath + "/p2b"
	trxSYSPath  = trxBasepath + "/sys"
)
