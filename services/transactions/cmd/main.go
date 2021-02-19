package main

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	serverPb "github.com/DaniilOr/microtracing/services/transactions/pkg/server"
	"github.com/DaniilOr/microtracing/services/transactions/pkg/transactions"
	"github.com/DaniilOr/microtracing/services/transactions/cmd/app"
)

const (
	defaultPort = "9999"
	defaultHost = "0.0.0.0"
	defaultDSN  = "postgres://app:pass@transactionsdb:5432/db"
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

	dsn, ok := os.LookupEnv("APP_DSN")
	if !ok {
		dsn = defaultDSN
	}

	if err := execute(net.JoinHostPort(host, port), dsn); err != nil {
		os.Exit(1)
	}
}

func execute(addr string, dsn string) error {
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

	grpcServer := grpc.NewServer()
	transactionsSVC := transactions.NewService(pool)
	server := app.NewServer(transactionsSVC, ctx)
	serverPb.RegisterTransactionsServerServer(grpcServer, server)
	return grpcServer.Serve(listener)
}