package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_generateQuery(t *testing.T) {
	t.Run("all fields provided", func(t *testing.T) {
		cc := "Amex"
		merchant := "ABC"
		category := "grocery"
		from := "2021/11/08"
		to := "2021/11/09"
		q, err := generateQuery(cc, merchant, category, from, to)
		require.NoError(t, err)
		require.Equal(t, "SELECT credit_card, amount, merchant, category, date FROM expenses WHERE credit_card = 'Amex' AND merchant = 'ABC' AND category = 'grocery' AND date >= '2021/11/08' AND date < '2021/11/09'", q)
	})
	t.Run("one fields provided", func(t *testing.T) {
		cc := "Amex"
		merchant := ""
		category := ""
		from := ""
		to := ""
		q, err := generateQuery(cc, merchant, category, from, to)
		require.NoError(t, err)
		require.Equal(t, "SELECT credit_card, amount, merchant, category, date FROM expenses WHERE credit_card = 'Amex'", q)
	})

	t.Run("some fields provided", func(t *testing.T) {
		cc := "Amex"
		merchant := "ABC"
		category := ""
		from := ""
		to := ""
		q, err := generateQuery(cc, merchant, category, from, to)
		require.NoError(t, err)
		require.Equal(t, "SELECT credit_card, amount, merchant, category, date FROM expenses WHERE credit_card = 'Amex' AND merchant = 'ABC'", q)
	})

	t.Run("No fields provided throws error", func(t *testing.T) {
		cc := ""
		merchant := ""
		category := ""
		from := ""
		to := ""
		_, err := generateQuery(cc, merchant, category, from, to)
		require.Error(t, err)
	})

}
func Test_ExpenseStore(t *testing.T) {
	f, err := os.CreateTemp("", "test_es.db")
	es, err := InitExpenseStore(f.Name(), false)
	defer es.Close()
	require.NoError(t, err)
	require.NotNil(t, es)
	ts := []*Transanction{
		{
			CreditCard: "Amex",
			Amount:     28.00,
			Merchant:   "XYZ",
			Date:       "2020/11/07",
			Category:   "Grocery",
		},
		{
			CreditCard: "Citi",
			Amount:     38.00,
			Merchant:   "ABC",
			Date:       "2021/11/07",
			Category:   "Grocery",
		},
	}
	for i := 0; i < len(ts); i++ {
		err := es.AddExpense(ts[i])
		require.NoError(t, err)
	}

	got, err := es.GetTransactions("Amex", "", "", "", "")
	require.NoError(t, err)
	require.Equal(t, ts[0], got[0])

	got, err = es.GetTransactions("Citi", "", "", "", "")
	require.NoError(t, err)
	require.Equal(t, ts[1], got[0])

}
