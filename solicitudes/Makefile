.PHONY: install strat test

instal:
	go mod tidy
	docker compose up -d	

start:
	go run cmd/main.go

test:
	go test ./... -v