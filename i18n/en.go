package i18n

var lang_en = Language{
	Months: [12]string{
		"January",
		"February",
		"March",
		"April",
		"May",
		"June",
		"July",
		"August",
		"September",
		"October",
		"November",
		"December",
	},
	Dictionary: map[string]string{
		"monthley_report":    "Monthly Report",
		"section_summary":    "Summary",
		"legal_notice_title": "Legal Notice",
		"legal_notice": `This report provides a record of transactions involving the SAVVA crypto token. You should carefully note the following:

* **Volatility of Crypto Assets:** The price of SAVVA, like all crypto assets, is highly volatile and subject to unpredictable fluctuations.

* **Informational Purposes Only:** This report is produced solely for the purpose of informing you about your account's activity.

* **No Financial Advice:** This report does *not* constitute financial advice. It should not be relied upon as the basis for any financial decisions.

* **Consult a Financial Specialist:** Before using this report in any way for financial purposes, you are strongly advised to consult with a qualified financial specialist.

* **Currency Conversion Disclaimer:** Any USD or EUR values presented in this report are based on the prevailing SAVVA token price at the time the report was generated. These values will change, potentially significantly, as the price of SAVVA fluctuates. These converted values should be treated with extreme caution.

* **Price Fluctuations:** The price of the SAVVA token can and will change unpredictably, and past performance is not indicative of future results.
`,
		"description":                   "Description",
		"summary_introduction":          "This report provides a summary of your SAVVA account activity from *%s* to *%s*.",
		"summary.savva_in":              "Deposited to the account",
		"summary.savva_out":             "Sent from the account",
		"summary.donations_contribute":  "Donations Contributed",
		"summary.donations_received":    "Donations Received",
		"summary.fund_contributed":      "Post Funds Contributed",
		"summary.fund_prizes_won":       "Post Funds Prizes Won",
		"summary.staking_in":            "Staking Deposited",
		"summary.staking_out":           "Staking Withdrawn",
		"summary.staking_staked":        "Added to staking",
		"summary.club_buy":              "Spent on sponsoring authors",
		"summary.club_claimed":          "Receved from sponsors",
		"summary.fundrase_contributed":  "Fundraise Contributed",
		"summary.fundrase_received":     "Fundraise Received",
		"summary.paid_for_promotion":    "Paid for promotion",
		"summary.nft_share_received":    "NFT Share from Post Funds",
		"summary.nft_sold_received":     "NFT Sold Received",
		"summary.nft_auctions_bids":     "NFT Auctions Bids",
		"summary.nft_auctions_received": "NFT Auctions Received",
		"section_my_authors":            "My Authors",
		"my_authors_introduction":       "These are the SAVVA users you support. The weekly payment amounts reflect the values at the time this report was generated. Your total weekly support is %s.",
		"account":                       "Account",
		"total":                         "Total",
		"my_share":                      "My Share",
	},
}

// Add more translations as needed
