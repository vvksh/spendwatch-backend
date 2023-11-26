package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var es *ExpenseStore

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func main() {
	var err error
	connectionString := os.Getenv("MYSQL_STRING")
	if connectionString == "" {
		log.Fatal("no connection string declared in env")
	}
	fmt.Println("found connection string", connectionString)

	es, err := InitExpenseStore(connectionString)
	if err != nil {
		log.Panic(err)
	}
	defer es.Close()

	fmt.Printf("Listing on 8000")
	m := mux.NewRouter()
	m.HandleFunc("/expenses", getExpensesSummary).Methods(http.MethodGet)
	http.ListenAndServe(":8000", m)
}

func getExpensesSummary(w http.ResponseWriter, r *http.Request) {
	groupBy := r.URL.Query().Get("groupBy")
	var out any
	var err error
	var httpErrorType int
	if groupBy == "" {
		err = fmt.Errorf("Invalid groupby %s", groupBy)
		httpErrorType = http.StatusBadRequest
	} else if groupBy == "month" {
		out, err = es.GetMonthlyExpensesSummary()
		httpErrorType = http.StatusInternalServerError
	} else if groupBy == "credit_card" {
		out, err = es.GetCreditCardExpensesSummary()
		httpErrorType = http.StatusInternalServerError
	} else {
		err = fmt.Errorf("Invalid groupby %s", groupBy)
		httpErrorType = http.StatusBadRequest
	}

	if err != nil {
		w.WriteHeader(httpErrorType)
		w.Write([]byte(err.Error()))
		return
	}

	jsonOut, err := json.Marshal(out)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	enableCors(&w)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonOut)

}
