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
	docker run -d -v todo:/app -P go_final_project:v1.1.0

testall:
	go test ./tests/...