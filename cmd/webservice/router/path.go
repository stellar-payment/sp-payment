package router

const (
	basePath = "/payment/api/v1"
	PingPath = basePath + "/ping"

	// ----- Customers
	customerBasepath = basePath + "/customers"
	customerMePath   = customerBasepath + "/me"
	customerIDPath   = customerBasepath + "/:customerID"

	// ----- Merchants
	merchantBasepath = basePath + "/merchants"
	merchantMePath   = merchantBasepath + "/me"
	merchantIDPath   = merchantBasepath + "/:merchantID"

	// ----- Accounts
	accountBasepath = basePath + "/accounts"
	accountMePath   = accountBasepath + "/me"
	accountIDPath   = accountBasepath + "/:accountID"

	// ----- Transactions
	trxBasepath = basePath + "/transactions"
	trxIDPath   = trxBasepath + "/:trxID"
	trxP2PPath  = trxBasepath + "/p2p"
	trxP2BPath  = trxBasepath + "/p2b"
	trxSYSPath  = trxBasepath + "/sys"

	// ----- Settlements
	settlementBasepath = basePath + "/settlements"
	settlementIDPath   = settlementBasepath + "/:settlementID"

	// ----- Beneficiaries
	beneficiaryBasepath    = basePath + "/beneficiaries"
	beneficiaryIDPath      = beneficiaryBasepath + "/:beneficiaryID"
	beneficiaryPreviewPath = beneficiaryBasepath + "/preview"
)
