fixtures:
	go build -o fixtures cmd/fixtures/main.go

build:
	env GOOS=linux GOARCH=amd64 go build