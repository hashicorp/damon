build:
	go build -o bin/damon ./cmd/damon

run:
	./bin/damon

test:
	go test ./...
