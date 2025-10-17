# Variables
BINARY_NAME=documentos

# Colores para Windows (sin códigos ANSI para evitar problemas)
GREEN=[92m
YELLOW=[93m
WHITE=[97m
RED=[91m
RESET=[0m

.PHONY: help install start

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
