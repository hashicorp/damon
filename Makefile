build:
	go build -o bin/damon ./cmd/damon

run:
	./bin/damon

install-osx:
	cp ./bin/damon /usr/local/bin/damon

test:
	go test ./...
