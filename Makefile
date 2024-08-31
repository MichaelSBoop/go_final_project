ifneq (,$(wildcard ./.env*))
	include .env*
	export
endif

run:
	go run .

exe:
	go build -o ./app .

image: 
	@echo "building image..."
	docker build --tag go_final_project:v1.1.0 .
	
container: image
	@echo "creating container..."
	docker volume create todo
	docker run -d -v todo:$(TODO_DBFILE) -p 7540:7540 go_final_project:v1.1.0

testall: container
	go test ./tests/...