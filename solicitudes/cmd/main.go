// Para compilar: go run cmd/main.go

package main

import (
	"fmt"
	"log"
	"os"

    "github.com/kramirez/solicitudes/internal/solicitud"
    "github.com/kramirez/solicitudes/pkg/bootstrap"
    "github.com/kramirez/solicitudes/pkg/handler"
    "github.com/kramirez/solicitudes/pkg/httpclient"
)

func main() {
	//Iniciar Logger
	logger := bootstrap.InitLogger()

	//Cargar variables de entorno
	bootstrap.InitEnv()

	//Conectar a la base de datos
	db, err := bootstrap.DBConnection()
	if err != nil {
		log.Fatal("Error al conectar a la base de datos", err)
	}

	logger.Println("Base de datos conectada exitosamente")

	// Crear cliente para el microservicio de documentos
	documentoClient := httpclient.NewDocumentoClient("http://localhost:8083")

	// Inicializar repositorio
	solicitudRepo := solicitud.NewRepository(db)

	// Inicializar servicio con el cliente de documentos
	service := solicitud.NewService(solicitudRepo, logger, documentoClient)

	// Inicializar endpoint
	endpoint := solicitud.NewEndpoint(service)

	//Configurar rutas
	router := handler.SetupRoutes(endpoint)

	//Obtener puerto del servicio
	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8082"
	}

	logger.Printf("Servicio de solicitudes iniciado en el puerto %s", port)

	//Iniciar el servidor
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal("Error al iniciar el servidor", err)
	}

}
