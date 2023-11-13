package main

import (
	"log"
	"time"
)

func fetchRecentExpenses(es *ExpenseStore, limit int) error {
	log.Printf("Fetching %d expenses \n", limit)
	// TODO Replace with calls to gmail
	ccs := []string{"Amex", "Citi", "Chase", "Discover"}
	merchants := []string{"ABC", "XYZ", "PQR", "BOK", "PLOP"}
	amounts := []float64{28.00, 76.00, 87.65, 900.98, 345.78, 654.76}
	categories := []string{"Restaurant", "groceries", "fitness"}
	dates := []string{"2020/11/08", "2020/11/09", "2023/11/10", "2020/11/11", "2023/11/12", "2020/11/13"}
	count := 0

	for _, cc := range ccs {
		for _, d := range dates {
			for _, m := range merchants {
				for _, c := range categories {
					for _, a := range amounts {
						if count > limit {
							break
						}
						t := &Transanction{
							CreditCard: cc,
							Category:   c,
							Amount:     a,
							Merchant:   m,
							Date:       d,
						}
						log.Printf("Adding expense: %+v", t)
						err := es.AddExpense(t)
						if err != nil {
							return err
						}
						count++

					}
				}
			}
		}
	}
	return nil
}

func backgroundExpenseFetcher(es *ExpenseStore, limit int) {
	for {
		time.Sleep(30 * time.Minute)
		fetchRecentExpenses(es, limit)
	}
}
