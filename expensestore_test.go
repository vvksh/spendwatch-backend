package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ExpenseStore(t *testing.T) {
	f, err := os.CreateTemp("", "test_es.db")
	es, err := InitExpenseStore(f.Name())
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

	got, err := es.GetTransactionsByCreditCard("Amex")
	require.NoError(t, err)
	require.Equal(t, ts[0], got[0])

	got, err = es.GetTransactionsByCreditCard("Citi")
	require.NoError(t, err)
	require.Equal(t, ts[1], got[0])

}
