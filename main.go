package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var db *sql.DB

func getExpenses(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func main() {
	db, err := InitExpenseStore("./expense.db")
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	fmt.Printf("Listing on 8000")
	m := mux.NewRouter()
	m.HandleFunc("/expenses", getExpenses).Methods(http.MethodGet)
	http.ListenAndServe(":8000", m)
}
