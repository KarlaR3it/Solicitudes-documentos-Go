package main

import (
	"log"

	"github.com/kramirez/usuarios-service/internal/usuario"
	"github.com/kramirez/usuarios-service/pkg/bootstrap"
	"github.com/kramirez/usuarios-service/pkg/handler"
)

func main() {
	// Inicializar logger
	logger := bootstrap.InitLogger()

	// Cargar variables de entorno
	bootstrap.InitEnv()

	// Conectar a la base de datos
	db, err := bootstrap.DBConnection()
	if err != nil {
		log.Fatal("Error al conectar a la base de datos:", err)
	}

	logger.Println("Base de datos conectada exitosamente")
	logger.Println("Servicio de usuarios iniciado")")

	// Mantener el servicio corriendo
	select {}
}