package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const dbPath = "./expense.db"

var db *ExpenseStore

func getExpenses(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func main() {
	es, err := InitExpenseStore(dbPath)
	go backgroundExpenseFetcher(es, 5)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	fmt.Printf("Listing on 8000")
	m := mux.NewRouter()
	m.HandleFunc("/expenses", getExpenses).Methods(http.MethodGet)
	http.ListenAndServe(":8000", m)
}
