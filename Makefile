fixtures:
	go build -o fixtures cmd/fixtures/main.go

build:
	env GOOS=linux GOARCH=arm64 go build