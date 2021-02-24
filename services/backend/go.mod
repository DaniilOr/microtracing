module github.com/DaniilOr/microtracing/services/backend

go 1.15

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	github.com/DaniilOr/microtracing/services/auth v0.0.0-20210224132015-35aa5ba4eda7
	github.com/DaniilOr/microtracing/services/transactions v0.0.0-20210218143633-fd506a6b2876
	github.com/go-chi/chi v1.5.2
	go.opencensus.io v0.22.6
	google.golang.org/grpc v1.35.0
)
