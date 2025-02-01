module github.com/mendes113/protosocket

go 1.22.0

toolchain go1.23.5

require (
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.3
	go.opentelemetry.io/otel/metric v1.34.0
	go.opentelemetry.io/otel/trace v1.34.0
	go.uber.org/zap v1.27.0
	golang.org/x/time v0.9.0
	google.golang.org/protobuf v1.36.4
)

require go.uber.org/multierr v1.10.0 // indirect

require (
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel v1.34.0
)
