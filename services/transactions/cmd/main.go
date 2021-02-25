package main

import (
	"context"
	"contrib.go.opencensus.io/exporter/jaeger"
	"github.com/DaniilOr/microtracing/services/transactions/cmd/app"
	serverPb "github.com/DaniilOr/microtracing/services/transactions/pkg/server"
	"github.com/DaniilOr/microtracing/services/transactions/pkg/transactions"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"os"
)

const (
	defaultPort = "8888"
	defaultHost = "0.0.0.0"
	defaultDSN  = "postgres://app:pass@transactionsdb:5432/db"
	defaultCertificatePath = "put path here"
	 defaultPrivateKeyPath = "put path here"
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

	dsn, ok := os.LookupEnv("APP_DSN")
	if !ok {
		dsn = defaultDSN
	}
	err := InitJaeger("transactions")
	if err != nil{
		log.Println(err)
	}
	certificatePath, ok := os.LookupEnv("APP_CERT_PATH")
	if !ok {
		certificatePath = defaultCertificatePath
	}

	privateKeyPath, ok := os.LookupEnv("APP_PRIVATE_KEY_PATH")
	if !ok {
		privateKeyPath = defaultPrivateKeyPath
	}
	if err := execute(net.JoinHostPort(host, port), certificatePath, privateKeyPath, dsn); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func execute(addr string, certificate string, key string, dsn string) error {
	creds, err := credentials.NewServerTLSFromFile(certificate, key)
	if err != nil {
		log.Fatalf("Failed to setup TLS: %v", err)
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		log.Print(err)
		return err
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	transactionsSVC := transactions.NewService(pool)
	server := app.NewServer(transactionsSVC, ctx)
	serverPb.RegisterTransactionsServerServer(grpcServer, server)
	return grpcServer.Serve(listener)
}