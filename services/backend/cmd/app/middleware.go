package app

import (
	"context"
	"errors"
	"net/http"
	authorization "github.com/DaniilOr/microtracing/services/auth/pkg/auth"
	"github.com/DaniilOr/microtracing/services/backend/pkg/jwt/symmetric"
)

var ErrNoAuth = errors.New("no auth in context")

type AuthFunc func(ctx context.Context, token string) (userID int64, err error)

var authContextKey = &contextKey{"auth context"}

type contextKey struct {
	name string
}
type Data struct {
	UserID int64    `json:"userId"`
	Roles  []string `json:"roles"`
	Issued int64    `json:"iat"`
	Expire int64    `json:"exp"`
}
func (c *contextKey) String() string {
	return c.name
}

func Auth(authFunc AuthFunc) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			token := request.Header.Get("Authorization")
			if token == "" {
				writer.WriteHeader(http.StatusUnauthorized)
				return
			}

			auth, err := authFunc(request.Context(), token)
			if err != nil {
				// упрощённый вариант, нужно ещё добавить проверку на то, что удалённый сервис "отвалился"
				if !errors.As(err, &authorization.ErrUserNotFound){
					key := []byte("some secter key goes here")
					verified, err := symmetric.Verify(token, key)
					if err != nil{
						writer.WriteHeader(http.StatusForbidden)
						return
					}
					if !verified{
						writer.WriteHeader(http.StatusUnauthorized)
						return
					}
					var decoded *Data
					err = symmetric.Decode(token, &decoded)
					if err != nil {
						writer.WriteHeader(http.StatusForbidden)
						return
					}
					ctx := context.WithValue(request.Context(), authContextKey, decoded.UserID)
					request = request.WithContext(ctx)
					handler.ServeHTTP(writer, request)
					return
				}
				writer.WriteHeader(http.StatusForbidden)
				return
			}

			ctx := context.WithValue(request.Context(), authContextKey, auth)
			request = request.WithContext(ctx)
			handler.ServeHTTP(writer, request)
		})
	}
}

func AuthFrom(ctx context.Context) (int64, error) {
	if value := ctx.Value(authContextKey); value != nil {
		if id, ok := value.(int64); ok {
			return id, nil
		}
	}
	return 0, ErrNoAuth
}
