package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

// Interacts with database
type ExpenseStore struct {
	db *sql.DB
}

func InitExpenseStore(dbPath string) (*ExpenseStore, error) {
	db, err := sql.Open("mysql", dbPath)
	if err != nil {
		return nil, err
	}
	if db == nil {
		return nil, fmt.Errorf("no db instantiated")
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

func (e *ExpenseStore) GetMonthlyExpensesSummary() (map[string]string, error) {
	output := map[string]string{}
	// get last 3 months of data
	queryStmt := "select MONTHNAME(date) as month, sum(amount) as sum from expenses where date >= CURDATE()- INTERVAL 3 MONTH group by MONTHNAME(date);"

	rows, err := e.db.Query(queryStmt)
	if err != nil {
		return output, err
	}
	defer rows.Close()
	for rows.Next() {
		var month, sum string
		err = rows.Scan(&month, &sum)
		output[month] = sum
	}
	return output, nil
}

func (e *ExpenseStore) GetCreditCardExpensesSummary() (map[string]string, error) {
	output := map[string]string{}
	// get last 3 months of data
	queryStmt := "select credit_card, sum(amount) from expenses as sum where date >= CURDATE()- INTERVAL 3 MONTH group by credit_card;"

	rows, err := e.db.Query(queryStmt)
	if err != nil {
		return output, err
	}
	defer rows.Close()
	for rows.Next() {
		var credit_card, sum string
		err = rows.Scan(&credit_card, &sum)
		output[credit_card] = sum

	}
	return output, nil
}

func (e *ExpenseStore) Close() {
	log.Printf("Closing db \n")
	e.db.Close()
}
