package main

type Transanction struct {
	CreditCard string  `json:"credit_card"`
	Amount     float64 `json:"amount"`
	Merchant   string  `json:"merchant"`
	Date       string  `json:"date"`
	Category   string  `json:"category"`
}
