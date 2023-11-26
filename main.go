package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const dbPath = "./expense.db"

var es *ExpenseStore

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func getExpenses(w http.ResponseWriter, r *http.Request) {
	cc := r.URL.Query().Get("cc")
	category := r.URL.Query().Get("category")
	merchant := r.URL.Query().Get("merchant")
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	ts, err := es.GetTransactions(cc, merchant, category, from, to)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	jsonOut, err := json.Marshal(ts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	enableCors(&w)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonOut)
}

func main() {
	var err error
	es, err = InitExpenseStore(dbPath, true)
	go backgroundExpenseFetcher(es, 5)
	if err != nil {
		log.Panic(err)
	}
	defer es.Close()

	fmt.Printf("Listing on 8000")
	m := mux.NewRouter()
	m.HandleFunc("/expenses", getExpenses).Methods(http.MethodGet)
	http.ListenAndServe(":8000", m)
}
