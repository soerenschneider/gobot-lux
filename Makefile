build-raspberry:
	GOARM=6 GOARCH=arm GOOS=linux go build -o weatherbot-armv6 cmd/main.go

unittest:
	go test ./...
