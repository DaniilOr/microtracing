package app

import (
	"context"
	serverPb "github.com/DaniilOr/microtracing/services/transactions-api/pkg/server"
	"github.com/DaniilOr/microtracing/services/transactions-api/pkg/transactions"
	"log"
)

type Server struct {
	transactionsSvc *transactions.Service
	ctx context.Context
}

func NewServer(transactionsSvc *transactions.Service, ctx context.Context) *Server {
	return &Server{transactionsSvc: transactionsSvc, ctx: ctx}
}

func (s *Server) Transactions(ctx context.Context, request * serverPb.TransactionsRequest) (*serverPb.TransactionsResponse, error){

	userID := request.UserID

	records, err := s.transactionsSvc.Transactions(ctx, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var response serverPb.TransactionsResponse
	for _, record := range records {
		response.Items = append(response.Items, &serverPb.Transaction{
			Id:       record.ID,
			UserId:   record.UserID,
			Category: record.Category,
			Amount:   record.Amount,
			Created:  record.Created,
		})
	}
	return &response, nil
}