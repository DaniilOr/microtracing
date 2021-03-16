package main

import (
	"contrib.go.opencensus.io/exporter/jaeger"
	"github.com/DaniilOr/microtracing/services/backend/cmd/app"
	"github.com/DaniilOr/microtracing/services/backend/pkg/auth"
	"github.com/DaniilOr/microtracing/services/backend/pkg/transactions"
	"github.com/go-chi/chi"
	"go.opencensus.io/trace"
	"log"
	"net"
	"net/http"
	"os"
)


const (
	defaultPort               = "9999"
	defaultHost               = "0.0.0.0"
	defaultAuthURL            = "netology.local:8080"
	defaultTransactionsAPIURL = "netology.local:8888"
)
func InitJaeger(serviceName string) error{
	exporter, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint: "jaeger:6831",
		Process: jaeger.Process{
			ServiceName: serviceName,
			Tags: []jaeger.Tag{
				jaeger.StringTag("hostname", "localhost"),
			},
		},
	})
	if err != nil {
		return err
	}
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.AlwaysSample(),
	})
	return nil
}
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
	err := InitJaeger("backed")
	if err != nil{
		log.Println(err)
		os.Exit(1)
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
	mux := chi.NewRouter()
	transactionsSVC, err := transactions.Init(transactionsAPIURL)
	if err != nil{
		return err
	}
	application := app.NewServer(authSvc, transactionsSVC, mux)
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