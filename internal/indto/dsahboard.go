package indto

type GenericDashboardGraph struct {
	Key   any     `db:"key"`
	Value float64 `db:"value"`
}

type TransactionMetaDashboard struct {
	SenderName    string  `db:"sender_name"`
	RecipientName string  `db:"recipient_name"`
	Nominal       float64 `db:"nominal"`
	TrxDate       string  `db:"trx_date"`
}

type AdminDashboard struct {
	PeerTrxCount     int64                   `db:"peer_trx_count"`
	MerchantTrxCount int64                   `db:"merchant_trx_count"`
	SystemTrxCount   int64                   `db:"system_trx_count"`
	TotalCustomers   int64                   `db:"total_customers"`
	TotalMerchants   int64                   `db:"total_merchants"`
	TrxTraffic       []GenericDashboardGraph `db:"-"`
}

type MerchantDashboardParams struct {
	AccountID  string
	MerchantID string
}

type MerchantDashboard struct {
	TrxCount           int64   `db:"trx_count"`
	TrxNominal         float64 `db:"trx_nominal"`
	SettlementNominal  float64 `db:"settlement_nominal"`
	BeneficiaryNominal float64 `db:"beneficiary_nominal"`
}

type CustomerDashboardParams struct {
	AccountID  string
	CustomerID string
}

type CustomerDashboard struct {
	PeerTrxCount       int64   `db:"peer_trx_count"`
	PeerTrxNominal     float64 `db:"peer_trx_nominal"`
	MerchantTrxCount   int64   `db:"merchant_trx_count"`
	MerchantTrxNominal float64 `db:"merchant_trx_nominal"`
}
