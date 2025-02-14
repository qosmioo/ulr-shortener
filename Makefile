.PHONY: run

run:
	@docker-compose up -d

down:
	@docker-compose down

tests:
	@go test ./...