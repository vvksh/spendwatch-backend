package main

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var es *ExpenseStore
var password string
var logger zap.SugaredLogger

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func main() {
	var err error
	logger := NewLogger()
	if err != nil {
		log.Fatal(err)
	}
	password = os.Getenv("PASSWORD")
	if password == "" {
		log.Fatal("no connection string declared in env")
	}
	connectionString := os.Getenv("MYSQL_STRING")
	if connectionString == "" {
		log.Fatal("no connection string declared in env")
	}
	fmt.Println("found connection string: %s", connectionString)
	port := os.Getenv("PORT")
	if connectionString == "" {
		log.Fatal("no port in env")
	}
	logger.Infof("found port", port)
	addr := fmt.Sprintf(":%s", port)
	es, err = InitExpenseStore(connectionString)
	if err != nil {
		log.Panic(err)
	}
	defer es.Close()

	logger.Infof("Listing on %s", port)
	m := mux.NewRouter()
	m.HandleFunc("/expenses", rateLimiter(getExpensesSummary)).Methods(http.MethodGet)
	http.ListenAndServe(addr, m)
}

func getExpensesSummary(w http.ResponseWriter, r *http.Request) {
	logger.Infof("Handling request wirh query: %s \n", r.URL.RawQuery)
	enableCors(&w)
	providedPassword := r.URL.Query().Get("pwd")
	groupBy := r.URL.Query().Get("groupBy")

	if providedPassword != password {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Password incorrect"))
		return
	}

	var out any
	var err error
	var httpErrorType int
	if groupBy == "" {
		err = fmt.Errorf("No groupBy parameter passed")
		httpErrorType = http.StatusBadRequest
	} else if groupBy == "month" {
		out, err = es.GetMonthlyExpensesSummary()
		httpErrorType = http.StatusInternalServerError
	} else if groupBy == "credit_card" {
		out, err = es.GetCreditCardExpensesSummary()
		httpErrorType = http.StatusInternalServerError
	} else if groupBy == "category" {
		out, err = es.GetCategoryExpensesSummary()
		httpErrorType = http.StatusInternalServerError
	} else {
		err = fmt.Errorf("Invalid groupby %s", groupBy)
		httpErrorType = http.StatusBadRequest
	}

	if err != nil {
		log.Printf(err.Error())
		w.WriteHeader(httpErrorType)
		w.Write([]byte(err.Error()))
		return
	}

	jsonOut, err := json.Marshal(out)
	if err != nil {
		log.Printf(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonOut)
}

func NewLogger() *zap.SugaredLogger {
	cfg := zap.NewProductionConfig()
	LOG_FILE := os.Getenv("LOG_FILE")
	if LOG_FILE == "" {
		LOG_FILE = "/var/log/app/spendwatch_backend.log"
	}
	cfg.OutputPaths = []string{
		LOG_FILE,
	}
	l, err := cfg.Build()
	if err != nil {
		log.Fatal(err)
	}
	return l.Sugar()
}
