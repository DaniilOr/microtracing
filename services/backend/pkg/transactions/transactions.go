package transactions

import (
	"context"
	"encoding/json"
	serverPb "github.com/DaniilOr/microtracing/services/transactions/pkg/server"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
)
type Service struct{
	client serverPb.TransactionsServerClient
}

type ResponseDTO struct{
	Category string `json:"category"`
	Cost int64 `json:"amount"`
}
func Init(addr string) (*Service, error){
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return  nil, err
	}
	client :=  serverPb.NewTransactionsServerClient(conn)
	server := Service{client: client}
	return &server, nil
}

func (s*Service) Transactions(ctx context.Context, userId int64) (data []byte, err error) {
	ctx, span := trace.StartSpan(context.Background(), "route: transactions")
	defer span.End()
	response, err := s.client.Transactions(ctx, &serverPb.TransactionsRequest{UserID: userId})
	if err != nil{
		return nil, err
	}
	resp := make([]ResponseDTO, len(response.Items))
	for i, trans := range response.Items{
		resp[i] =  ResponseDTO{Category: trans.Category, Cost: trans.Amount}
	}
	data, err = json.Marshal(resp)
	if err != nil{
		return nil, err
	}
	return data, nil
}
