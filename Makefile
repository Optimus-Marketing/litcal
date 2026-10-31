.PHONY: test test-all lint cover golden

test:
	go test -race -short ./...

test-all:
	go test ./...

lint:
	test -z "$$(gofmt -l .)"
	go vet ./...
	staticcheck ./...

cover:
	go test -short -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

golden:
	go test -run TestYearMatchesGolden -update .
