package data

import (
	"database/sql"
	"math/big"
	"slices"

	"github.com/AlexNa-Holdings/savva-reports/cmn"
)

type DomainRecord struct {
	Domain       string
	Amount       *big.Int
	CurrentFrame int
	TilFrame     int
}

type Sponsored struct {
	Author       string
	Domain       string
	Member       string
	Domains      []DomainRecord
	TotalAmount  *big.Int
	TotalFromAll *big.Int
}

// CREATE TABLE public.clubs (
// 	"domain" varchar(255) NOT NULL,
// 	author_addr varchar(42) NOT NULL,
// 	current_frame int4 DEFAULT 0 NULL,
// 	CONSTRAINT clubs_pkey PRIMARY KEY (domain, author_addr)
// );

// CREATE TABLE public.clubs_members (
// 	"domain" varchar(255) NOT NULL,
// 	author_addr varchar(42) NOT NULL,
// 	member_addr varchar(42) NOT NULL,
// 	amount numeric(78) DEFAULT 0 NULL,
// 	til_frame int4 DEFAULT 0 NULL,
// 	CONSTRAINT clubs_members_pkey PRIMARY KEY (domain, author_addr, member_addr)
// );

func GetSponsoredBy(address string) ([]Sponsored, error) {
	rows, err := cmn.C.DB.Query(`
SELECT 
	cm.domain, 
	cm.author_addr, 
	cm.amount, 
	cm.til_frame,
	c.current_frame,
	author_totals.total_author_amount
FROM clubs_members cm
JOIN clubs c 
	ON cm.domain = c.domain AND cm.author_addr = c.author_addr
LEFT JOIN (
	SELECT author_addr, SUM(amount) AS total_author_amount
	FROM clubs_members
	GROUP BY author_addr
) AS author_totals
	ON cm.author_addr = author_totals.author_addr
WHERE cm.member_addr = $1
ORDER BY cm.author_addr, cm.domain;
	`, address)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Sponsored
	var current Sponsored

	for rows.Next() {
		var domain string
		var author string
		var amount sql.NullString
		var tilFrame sql.NullInt32
		var currentFrame sql.NullInt32
		var totalAuthorAmount sql.NullString

		err := rows.Scan(&domain, &author, &amount, &tilFrame, &currentFrame, &totalAuthorAmount)
		if err != nil {
			return nil, err
		}

		// Initialize new group if author changes
		if current.Author != author {
			if current.Author != "" {
				results = append(results, current)
			}
			current = Sponsored{
				Author:      author,
				Member:      address,
				TotalAmount: big.NewInt(0),
				Domains:     []DomainRecord{},
			}

			if totalAuthorAmount.Valid {
				current.TotalFromAll = new(big.Int)
				current.TotalFromAll.SetString(totalAuthorAmount.String, 10)
			} else {
				current.TotalFromAll = big.NewInt(0)
			}
		}

		amt := new(big.Int)
		if amount.Valid {
			amt.SetString(amount.String, 10)
		}
		current.TotalAmount.Add(current.TotalAmount, amt)

		current.Domain = domain // can store last one seen

		current.Domains = append(current.Domains, DomainRecord{
			Domain:       domain,
			Amount:       amt,
			CurrentFrame: int(currentFrame.Int32),
			TilFrame:     int(tilFrame.Int32),
		})
	}

	if current.Author != "" {
		results = append(results, current)
	}

	slices.SortFunc(results, func(a, b Sponsored) int {
		return b.TotalAmount.Cmp(a.TotalAmount)
	})

	return results, nil
}
