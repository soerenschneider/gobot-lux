build-raspberry:
	GOARM=6 GOARCH=arm GOOS=linux go build -o gobot-brightness-armv6 cmd/main.go

unittest:
	go test ./...
