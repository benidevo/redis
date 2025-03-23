test:
	codecrafters test

go-test:
	go test ./...g

run:
	sh scripts/run.sh

fmt:
	go fmt ./...
