package main

import (
	"encoding/json"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
	"os"
	"time"
)

var es *ExpenseStore
var password string
var logger *zap.SugaredLogger

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func main() {
	var err error
	err = sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	})
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)

	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	logger = NewLogger()
	if err != nil {
		log.Fatal(err)
	}
	password = os.Getenv("PASSWORD")
	if password == "" {
		log.Fatal("no password string declared in env")
	}
	connectionString := os.Getenv("MYSQL_STRING")
	if connectionString == "" {
		log.Fatal("no connection string declared in env")
	}
	log.Println("found connection string: %s", connectionString)
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
		"stdout",
	}
	cfg.ErrorOutputPaths = []string{
		LOG_FILE,
		"stderr",
	}
	cfgEncoder := zap.NewProductionEncoderConfig()
	cfgEncoder.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig = cfgEncoder
	logger := zap.Must(cfg.Build())
	return logger.Sugar()
}
