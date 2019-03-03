set GOOS=linux
set GOARCH=amd64
go build -ldflags="-s -w" main.go