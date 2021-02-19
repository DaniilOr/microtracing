package main

import (
	"github.com/DaniilOr/microtracing/services/backend/cmd/app"
	"github.com/DaniilOr/microtracing/services/backend/pkg/auth"
	"github.com/go-chi/chi"
	"log"
	"net"
	"net/http"
	"os"
)


const (
	defaultPort               = "9999"
	defaultHost               = "0.0.0.0"
	defaultAuthURL            = "auth:8080"
	defaultTransactionsAPIURL = "transactions:8888"
)

func main() {
	port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("APP_HOST")
	if !ok {
		host = defaultHost
	}

	authURL, ok := os.LookupEnv("APP_AUTH_URL")
	if !ok {
		authURL = defaultAuthURL
	}

	transactionsAPIURL, ok := os.LookupEnv("APP_TRANSACTIONS_URL")
	if !ok {
		transactionsAPIURL = defaultTransactionsAPIURL
	}

	if err := execute(net.JoinHostPort(host, port), authURL, transactionsAPIURL); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func execute(addr string, authURL string, transactionsAPIURL string) error {
	authSvc, err := auth.Init(authURL)
	if err != nil{
		return err
	}
	//transactionsSvc := transactions.NewService(transactionsAPIURL)

	mux := chi.NewRouter()

	application := app.NewServer(authSvc, nil, mux)
	err = application.Init()
	if err != nil {
		log.Print(err)
		return err
	}

	server := &http.Server{
		Addr:    addr,
		Handler: application,
	}
	return server.ListenAndServe()
}