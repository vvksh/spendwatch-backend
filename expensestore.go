package main

import (
	"database/sql"
	"fmt"
	"strings"

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
	create table  if not exists expenses (id integer not null primary key, credit_card text, amount real, merchant text, date text, category text);
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

func generateQuery(cc string, merchant string, category string, from string, to string) (string, error) {
	query := "SELECT credit_card, amount, merchant, category, date FROM expenses WHERE "
	conditionals := []string{}
	if cc != "" {
		conditionals = append(conditionals, fmt.Sprintf("credit_card = '%s'", cc))
	}
	if merchant != "" {
		conditionals = append(conditionals, fmt.Sprintf("merchant = '%s'", merchant))
	}
	if category != "" {
		conditionals = append(conditionals, fmt.Sprintf("category = '%s'", category))
	}
	if from != "" {
		conditionals = append(conditionals, fmt.Sprintf("date >= '%s'", from))
	}
	if to != "" {
		conditionals = append(conditionals, fmt.Sprintf("date < '%s'", to))
	}
	if len(conditionals) == 0 {
		return "", fmt.Errorf("No conditions passed")
	}
	query += strings.Join(conditionals, " AND ")
	return query, nil
}

func (e *ExpenseStore) GetTransactions(cc string, merchant string, category string, from string, to string) ([]Transanction, error) {
	output := []Transanction{}
	queryStmt, err := generateQuery(cc, merchant, category, from, to)
	if err != nil {
		return output, err
	}
	rows, err := e.db.Query(queryStmt)
	if err != nil {
		return output, err
	}
	defer rows.Close()
	for rows.Next() {
		var t Transanction
		err = rows.Scan(&t.CreditCard, &t.Amount, &t.Merchant, &t.Category, &t.Date)
		output = append(output, t)
	}
	return output, nil
}

func (e *ExpenseStore) AddExpense(t *Transanction) error {
	insertStmt := fmt.Sprintf("insert into expenses values(NULL,'%s', %f, '%s', '%s', '%s')", t.CreditCard, t.Amount, t.Merchant, t.Date, t.Category)
	_, err := e.db.Exec(insertStmt)
	return err
}

func (e *ExpenseStore) Close() {
	e.db.Close()
}
