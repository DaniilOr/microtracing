package main

import (
	"context"
	"contrib.go.opencensus.io/exporter/jaeger"
	"github.com/DaniilOr/microtracing/services/auth/cmd/app"
	"github.com/DaniilOr/microtracing/services/auth/pkg/auth"
	serverPb "github.com/DaniilOr/microtracing/services/auth/pkg/server"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"os"
)

const (
	defaultPort = "8080"
	defaultHost = "0.0.0.0"
	defaultDSN  = "postgres://app:pass@authdb:5432/db"
	defaultCertificatePath = "./tls/certificate.pem"
	defaultPrivateKeyPath = "./tls/key.pem"
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
	err := InitJaeger("auth")
	if err != nil{
		log.Println(err)
		os.Exit(1)
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
	creds, err := credentials.NewServerTLSFromFile(defaultCertificatePath, defaultPrivateKeyPath)
	grpcServer := grpc.NewServer(grpc.Creds(creds), grpc.StatsHandler(&ocgrpc.ServerHandler{}))
	authSVC := auth.NewService(pool)
	server := app.NewServer(authSVC, ctx)
	serverPb.RegisterAuthServerServer(grpcServer, server)
	return grpcServer.Serve(listener)
}

