.PHONY: install start test cover clean

# Instalacion y configuracion
install:
	go mod tidy
	docker compose up -d

# Ejecutar aplicacion
start:
	go run cmd/main.go

# Ejecutar tests
test:
	go test ./internal/solicitud -v

# Coverage de endpoints (>75% requerido)
cover: clean
	go test ./internal/solicitud -coverprofile=coverage.out -covermode=count
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out

# Limpiar archivos temporales
clean:
	@if exist coverage.out del coverage.out 2>nul
	@if exist coverage.html del coverage.html 2>nul