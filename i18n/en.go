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
		"description":                  "Description",
		"summary.savva_in":             "Deposited to the account",
		"summary.savva_out":            "Sent from the account",
		"summary.donations_contribute": "Donations Contributed",
		"summary.donations_received":   "Donations Received",
		"summary.fund_contributed":     "Fund Contributed",
		"summary.fund_prizes_won":      "Fund Prizes Won",
		"summary.staking_in":           "Staking Deposited",
		"summary.staking_out":          "Staking Withdrawn",
		"summary.staking_staked":       "Added to staking",
		"summary.club_buy":             "Spent on sponsoring authors",
		"summary.club_claimed":         "Receved from sponsors",
	},
}

// Add more translations as needed
