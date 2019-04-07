set GOOS=linux
set GOARCH=arm
go build -ldflags="-s -w"  main.go task.go store.go