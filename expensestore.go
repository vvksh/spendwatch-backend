package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

// Interacts with database
type ExpenseStore struct {
	db *sql.DB
}

func InitExpenseStore(dbPath string) (*ExpenseStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	e := &ExpenseStore{
		db: db,
	}
	return e, nil
}

func (e *ExpenseStore) AddExpense(t *Transanction) error {
	insertStmt := fmt.Sprintf("insert into expenses values(NULL,'%s', %f, '%s', '%s', '%s')", t.CreditCard, t.Amount, t.Merchant, t.Date, t.Category)
	_, err := e.db.Exec(insertStmt)
	return err
}

func (e *ExpenseStore) GetMonthlyExpensesSummary() ([]*MonthlySum, error) {
	output := []*MonthlySum{}
	// get last 3 months of data
	queryStmt := "select MONTHNAME(date) as month, sum(amount) as sum from expenses where date >= CURDATE()- INTERVAL 3 MONTH group by MONTHNAME(date);"

	rows, err := e.db.Query(queryStmt)
	if err != nil {
		return output, err
	}
	defer rows.Close()
	for rows.Next() {
		var m MonthlySum
		err = rows.Scan(&m.Month, &m.Sum)
		output = append(output, &m)
	}
	return output, nil
}

func (e *ExpenseStore) GetCreditCardExpensesSummary() ([]*CreditCardSum, error) {
	output := []*CreditCardSum{}
	// get last 3 months of data
	queryStmt := "select credit_card, sum(amount) from expenses as sum where date >= CURDATE()- INTERVAL 3 MONTH group by credit_card;"

	rows, err := e.db.Query(queryStmt)
	if err != nil {
		return output, err
	}
	defer rows.Close()
	for rows.Next() {
		var m CreditCardSum
		err = rows.Scan(&m.CreditCard, &m.Sum)
		output = append(output, &m)
	}
	return output, nil
}

func (e *ExpenseStore) Close() {
	e.db.Close()
}
