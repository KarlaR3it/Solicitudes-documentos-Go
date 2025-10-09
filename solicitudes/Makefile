# Variables
BINARY_NAME=solicitudes
TEST_PACKAGE=.\internal\solicitud\...
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

# Colores para Windows (sin códigos ANSI para evitar problemas)
GREEN=[92m
YELLOW=[93m
WHITE=[97m
RED=[91m
RESET=[0m

.PHONY: help install start test test-cover cover clean

## Help: Muestra esta ayuda
help:
	@echo.
	@echo Uso:
	@echo   make ^<target^>
	@echo.
	@echo Targets:
	@powershell -Command "Select-String -Pattern '^##.*$$' $(MAKEFILE_LIST) | ForEach-Object { $$_ -replace '##[ ]*','' } | ForEach-Object { $$parts = $$_ -split ':',2; '  {0,-15}{1}' -f $$parts[0], $$parts[1] }"

## Instalación: Instala dependencias y levanta los servicios
install:
	@echo Instalando dependencias...
	@go mod tidy
	@echo Iniciando servicios con Docker...
	@docker compose up -d

## Iniciar: Ejecuta la aplicación
start:
	@echo Iniciando la aplicación...
	@go run cmd/main.go

## Test: Ejecuta todas las pruebas
test:
	@echo Ejecutando pruebas...
	@go test -v $(TEST_PACKAGE)

## Test-cover: Ejecuta pruebas con cobertura y genera reporte HTML
test-cover: clean
	@echo Ejecutando pruebas con cobertura...
	@go test -v -coverprofile=$(COVERAGE_FILE) -covermode=atomic $(TEST_PACKAGE)
	@if exist $(COVERAGE_FILE) ( \
		echo Generando reporte de cobertura... & \
		go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML) & \
		echo Resumen de cobertura: & \
		go tool cover -func=$(COVERAGE_FILE) | findstr /C:"total:" \
	) else ( \
		echo Error: No se generaron datos de cobertura & \
		exit /b 1 \
	)

## Cover: Alias para test-cover
cover: test-cover

## Clean: Elimina archivos temporales
clean:
	@echo Limpiando archivos temporales...
	@if exist $(COVERAGE_FILE) del /q $(COVERAGE_FILE) 2>nul
	@if exist $(COVERAGE_HTML) del /q $(COVERAGE_HTML) 2>nul
	@echo Limpieza completada