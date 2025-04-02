package main

import (
	"database/sql"
	"io"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/AlexNa-Holdings/savva-reports/cmn"
	"github.com/AlexNa-Holdings/savva-reports/reports"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

func main() {

	type TestConfig struct {
		DBConnection string `yaml:"db_connection"`
		IPFSGateway  string `yaml:"ipfs_gateway"`
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).Level(zerolog.DebugLevel)

	yamlFile, err := os.ReadFile("/home/alexn/savva-reports.yaml")
	if err != nil {
		log.Error().Err(err).Msg("Failed to read config file")
		return
	}

	// Unmarshal the YAML file into the TestConfig struct
	var testConfig TestConfig
	err = yaml.Unmarshal(yamlFile, &testConfig)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal config file")
		return
	}

	IPFS_Gateway = testConfig.IPFSGateway
	cmn.C.SavvaTokenPrice = 0.0024470
	cmn.C.CurrencySymbol = "$"
	cmn.C.IPFS = Ipfs

	cmn.C.DB, err = sql.Open("postgres", testConfig.DBConnection)
	if err != nil {
		log.Fatal().Msgf("Cannot connect to DB. Error:(%v)", err)
		return
	}

	defer cmn.C.DB.Close()

	err = reports.Build("0xDf691828859e3Cb1e31E6D2F8A9b04F3B91A717f", 2025, 2, "AlexNaMonth.pdf", "en")
	//err = reports.Build("0x86002b3616cD8F8DC4C3cAC51571d833810B2718", 2025, 2, "IgorMonth.pdf", "en")
	//err = reports.Build("0xd20CEB10C3e90ba880c0a3824C9bcD1623F5D39A", 2025, 2, "AnelaMonth.pdf", "ru")

	if err != nil {
		log.Error().Err(err).Msg("Failed to build report")
	}
}

var IPFS_Gateway string

// Load file from the IPFS gateway
func Ipfs(cid string) []byte {

	url := IPFS_Gateway + cid

	resp, err := http.Get(url)
	if err != nil {
		log.Error().Err(err).Msgf("Error fetching IPFS file: %v", err)
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Error().Msgf("Error fetching IPFS file: %s", resp.Status)
		return nil
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Error reading IPFS file")
		return nil
	}
	return data
}
