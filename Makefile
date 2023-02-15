build:
	go build -o cloud-platform-go-get-module-bin .

run:
	go run cloud-platform-go-get-module-bin .

test:
	go test -race -covermode=atomic -coverprofile=c.out -v ./...

coverage:
	make test
	go tool cover -html=c.out

fmt:
	go fmt ./...
