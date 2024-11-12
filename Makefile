run:
	go run ./cmd/main.go

build:
	go build ./cmd/main.go

test:
	go test ./...

fuzz:
	go test -fuzz=Fuzz ./pkg/language

fmt:
	go fmt ./...

examples:
	go run ./cmd/examples/password/main.go
