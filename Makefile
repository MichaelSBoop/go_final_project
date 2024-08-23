include .env

run:
	go run .

exe:
	go build -o ./app/app .

docker-exe:
 	CGO_ENABLED=$(CGO_ENABLED) 
	GOOS=$(GOOS) 
	GOARCH=$(GOARCH) 
	go build -o ./app/app .

image: docker-exe
	@echo "building image..."
	docker build --build-arg TODO_PORT=${TODO_PORT} --build-arg TODO_DBFILE=${TODO_DBFILE} --tag go_final_project:v1.1.0 .
	
container: image
	@echo "creating container..."
	docker run -d --rm -p ${TODO_PORT}:${TODO_PORT} go_final_project:v1.1.0

testall: container
	go test ./tests/...