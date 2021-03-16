package auth

import (
	"context"
	serverPb "github.com/DaniilOr/microtracing/services/auth/pkg/server"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"os"
)
const defaultCertificatePath = "/tls/certificate.pem"

type Service struct{
	client serverPb.AuthServerClient
}

func Init(addr string) (*Service, error){
	certificatePath, ok := os.LookupEnv("APP_CERTIFICATE_PATH")
	if !ok {
		certificatePath = defaultCertificatePath
	}
	creds, err := credentials.NewClientTLSFromFile(certificatePath, "")
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}
	if err != nil {
		return  nil, err
	}
	client :=  serverPb.NewAuthServerClient(conn)
	server := Service{client: client}
	return &server, nil
}

func (s*Service) Token(ctx context.Context, login string, password string) (token string, err error) {
	ctx, span := trace.StartSpan(context.Background(), "route: token")
	defer span.End()
	response, err := s.client.Token(ctx, &serverPb.TokenRequest{Login: login, Password: password})
	if err != nil{
		return "", err
	}
	return response.Token, nil
}
func (s*Service) Id(ctx context.Context, token string) (int64, error) {
	response, err := s.client.Id(ctx, &serverPb.IdRequest{
		Token: token,
	})
	if err != nil{
		return 0, err
	}
	return response.UserId, nil
}
