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

	sqlStmt := `
	create table  if not exists expenses (id integer not null primary key, credit_card text, amount real, merchant text, date text);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, fmt.Errorf("%q: %s\n", err, sqlStmt)
	}
	e := &ExpenseStore{
		db: db,
	}
	fetchRecentExpenses(e, 100)
	return e, nil
}

func (e *ExpenseStore) GetTransactionsByCreditCard(creditCard string) ([]Transanction, error) {
	output := []Transanction{}
	rows, err := e.db.Query("select credit_card, amount, merchant, date from expenses where credit_card = ?", creditCard)
	if err != nil {
		return output, err
	}
	defer rows.Close()
	for rows.Next() {
		var t Transanction
		err = rows.Scan(&t.CreditCard, &t.Amount, &t.Merchant, &t.Date)
		output = append(output, t)
	}
	return output, nil
}

func (e *ExpenseStore) AddExpense(t *Transanction) error {
	_, err := e.db.Exec(fmt.Sprintf("insert into expenses values(NULL,'%s', %f, '%s', '%s')", t.CreditCard, t.Amount, t.Merchant, t.Date))
	return err
}

func (e *ExpenseStore) Close() {
	e.db.Close()
}
