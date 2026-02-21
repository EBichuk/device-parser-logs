lint: 
	golangci-lint run

run:
	docker-compose up -d

down:
	docker-compose down