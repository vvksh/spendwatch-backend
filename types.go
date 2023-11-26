package main

type Transanction struct {
	CreditCard string  `json:"credit_card"`
	Amount     float64 `json:"amount"`
	Merchant   string  `json:"merchant"`
	Date       string  `json:"date"`
	Category   string  `json:"category"`
}

type MonthlySum struct {
	Month string `json:"month"`
	Sum   string `json:"sum"`
}

type CreditCardSum struct {
	CreditCard string `json:"credit_card"`
	Sum        string `json:"sum"`
}
