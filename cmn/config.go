package cmn

import "database/sql"

type Config struct {
	DB              *sql.DB
	IPFS_Gateway    string
	SavvaTokenPrice float64
	CurrencySymbol  string
}

var (
	Margin     = 40.0
	PageWidth  = 595.28
	PageHeight = 841.89
)

var C *Config = &Config{}
