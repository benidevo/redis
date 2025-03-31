test:
	codecrafters test

go-test:
	go test ./...

run:
	sh scripts/run.sh

fmt:
	go fmt ./...
