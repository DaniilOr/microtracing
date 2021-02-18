package main

import (
	"context"
	"github.com/jackc/pgx/pgxpool"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
	"log"
	"microtracing/services/auth/cmd/app"
	"microtracing/services/auth/pkg/auth"
	serverPb "microtracing/services/auth/pkg/server"
	"net"
	"os"
)

const (
	defaultPort = "8080"
	defaultHost = "0.0.0.0"
	defaultDSN  = "postgres://app:pass@authdb:5432/db"
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

	grpcServer := grpc.NewServer(grpc.StatsHandler(&ocgrpc.ServerHandler{}))
	authSVC := auth.NewService(pool)
	server := app.NewServer(authSVC, ctx)
	serverPb.RegisterAuthServerServer(grpcServer, server)
	return grpcServer.Serve(listener)
}

