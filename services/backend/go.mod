module github.com/DaniilOr/microtracing/services/backend

go 1.15

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	github.com/DaniilOr/microtracing/services/auth v0.0.0-20210224153451-5a9ed5df142a
	github.com/DaniilOr/microtracing/services/transactions v0.0.0-20210224153451-5a9ed5df142a
	github.com/go-chi/chi v1.5.4
	github.com/jackc/fake v0.0.0-20150926172116-812a484cc733 // indirect
	go.opencensus.io v0.22.6
	google.golang.org/grpc v1.35.0
)
