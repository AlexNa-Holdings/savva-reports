package data

import (
	"database/sql"
	"log"
	"math/big"
	"time"

	"github.com/AlexNa-Holdings/savva-reports/cmn"
)

type HistoryRecord struct {
	Contract  sql.NullString
	Domain    sql.NullString
	Type      sql.NullString
	FromAddr  sql.NullString
	ToAddr    sql.NullString
	Amount    *big.Int
	Token     sql.NullString
	SavvaCid  sql.NullString
	Locales   sql.NullString
	Info      sql.NullString
	TxHash    sql.NullString
	TimeStamp time.Time
}

func GetHistory(address string, from, to *time.Time) ([]HistoryRecord, error) {
	rows, err := cmn.C.DB.Query(`SELECT contract, domain, type, from_addr, to_addr, amount, token, savva_cid, locales, info, tx_hash, time_stamp FROM history WHERE from_addr = $1 AND time_stamp BETWEEN $2 AND $3 ORDER BY time_stamp DESC`, address, from, to)
	if err != nil {
		log.Printf("Error querying history for %s: %v", address, err)
		return nil, err
	}
	defer rows.Close()

	var history []HistoryRecord
	for rows.Next() {
		var h HistoryRecord
		var n_amount sql.NullString

		err := rows.Scan(
			&h.Contract, &h.Domain, &h.Type, &h.FromAddr, &h.ToAddr, &n_amount,
			&h.Token, &h.SavvaCid, &h.Locales, &h.Info, &h.TxHash, &h.TimeStamp,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		if n_amount.Valid {
			h.Amount = new(big.Int)
			h.Amount.SetString(n_amount.String, 10)
		} else {
			h.Amount = big.NewInt(0)
		}

		history = append(history, h)
	}

	return history, nil
}
