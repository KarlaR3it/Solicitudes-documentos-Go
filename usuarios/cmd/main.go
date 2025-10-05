// Para compilar: go run cmd/main.go

package main

import (
	"fmt"
	"log"
	"os"

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

	// Inicializar capas
	repo := usuario.NewRepository(db)
	service := usuario.NewService(repo, logger)
	endpoint := usuario.NewEndpoint(service)

	// Configurar rutas
	router := handler.SetupRoutes(endpoint)

	// Obtener puerto del servicio
	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8081"
	}

	logger.Printf("Servicio de usuarios iniciado en el puerto %s", port)

	// Iniciar servidor
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}
