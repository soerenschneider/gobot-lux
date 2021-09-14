BINARY_NAME = gobot-brightness

build:
	go build -o $(BINARY_NAME) cmd/main.go

build-raspberry:
	GOARM=6 GOARCH=arm GOOS=linux go build -o $(BINARY_NAME)-armv6 cmd/main.go

coverage:
	go test ./... -covermode=count -coverprofile=coverage.out
	go tool cover -html=coverage.out -o=coverage.html
	go tool cover -func=coverage.out -o=coverage.out

unittest:
	go test ./...
