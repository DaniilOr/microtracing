package transactions

import (
	serverPb "github.com/DaniilOr/microtracing/services/transactions/pkg/server"
	"google.golang.org/grpc"
)
type Server struct{
	client serverPb.TransactionsServerClient
}

func (s*Server) NewServer(addr string)(*Server, error){
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return  nil, err
	}
	client :=  serverPb.NewTransactionsServerClient(conn)
	server := Service{client: client}
	return &server, nil
}
