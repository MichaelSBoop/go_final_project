run:
	go run .
build:
	go build -o ./cmd 

testall: run
	go test ./tests/...