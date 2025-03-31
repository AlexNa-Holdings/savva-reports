package cmn

import (
	"database/sql"

	"github.com/ethereum/go-ethereum/common"
)

// type SavvaContractsInfo struct {
// 	Config           ContractInfo
// 	SavvaToken       ContractInfo
// 	ContentNFT       ContractInfo
// 	ContentRegistry  ContractInfo
// 	Staking          ContractInfo
// 	ContentFund      ContractInfo
// 	NFTMarketplace   ContractInfo
// 	NFTAuction       ContractInfo
// 	SavvaFaucet      ContractInfo
// 	UserProfile      ContractInfo
// 	RandomOracle     ContractInfo
// 	Governance       ContractInfo
// 	ListMarket       ContractInfo
// 	Promo            ContractInfo
// 	BuyBurn          ContractInfo
// 	AuthorOfTheMonth ContractInfo
// 	AuthorsClubs     ContractInfo
// 	Fundraiser       ContractInfo
// 	SavvaNPO         ContractInfo
// 	SavvaNPOFactory  ContractInfo
// 	SavvaNPOEvents   ContractInfo
// }

type ContractAddresses struct {
	ContentNFT common.Address
}

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
