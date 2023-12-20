package dto

type GenericDashboardGraph struct {
	Key   any     `json:"key"`
	Value float64 `json:"value"`
}

type TransactionMetaDashboard struct {
	SenderName    string  `json:"sender_name"`
	RecipientName string  `json:"recipient_name"`
	Nominal       float64 `json:"nominal"`
	TrxDate       string  `json:"trx_date"`
	TrxType       int64   `json:"trx_type"`
}

type AdminDashboard struct {
	PeerTrxCount     int64                   `json:"peer_trx_count"`
	MerchantTrxCount int64                   `json:"merchant_trx_count"`
	SystemTrxCount   int64                   `json:"system_trx_count"`
	TotalCustomers   int64                   `json:"total_customers"`
	TotalMerchants   int64                   `json:"total_merchants"`
	TrxTraffic       []GenericDashboardGraph `json:"trx_traffic"`
}

type MerchantDashboard struct {
	AccountID          string                     `json:"account_id"`
	AccountBalance     float64                    `json:"account_balance"`
	TrxCount           int64                      `json:"trx_count"`
	TrxNominal         float64                    `json:"trx_nominal"`
	SettlementNominal  float64                    `json:"settlement_nominal"`
	BeneficiaryNominal float64                    `json:"beneficiary_nominal"`
	LastTrx            []TransactionMetaDashboard `json:"last_trx"`
}

type CustomerDashboard struct {
	AccountID          string                     `json:"account_id"`
	AccountBalance     float64                    `json:"account_balance"`
	PeerTrxCount       int64                      `json:"peer_trx_count"`
	PeerTrxNominal     float64                    `json:"peer_trx_nominal"`
	MerchantTrxCount   int64                      `json:"merchant_trx_count"`
	MerchantTrxNominal float64                    `json:"merchant_trx_nominal"`
	LastTrx            []TransactionMetaDashboard `json:"last_trx"`
}
