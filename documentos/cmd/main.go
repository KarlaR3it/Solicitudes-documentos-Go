// Para compilar: go run cmd/main.go

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kramirez/documentos/internal/documento"
	"github.com/kramirez/documentos/pkg/bootstrap"
	"github.com/kramirez/documentos/pkg/handler"
	"github.com/kramirez/documentos/pkg/httpclient"
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

	//Inicializar cliente HTTP para el servicio de solicitudes
	solicitudesServiceURL := os.Getenv("SOLICITUDES_SERVICE_URL")
	if solicitudesServiceURL == "" {
		solicitudesServiceURL = "http://localhost:8082"
	}
	solicitudesClient := httpclient.NewSolicitudClient(solicitudesServiceURL)
	logger.Printf("Cliente de solicitudes configurado: %s", solicitudesServiceURL)

	//Inicializar capas
	repo := documento.NewRepository(db)
	service := documento.NewService(repo, logger, solicitudesClient)
	endpoint := documento.NewEndpoint(service)

	//Configurar rutas
	router := handler.SetupRoutes(endpoint)

	//Obtener puerto del servicio
	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8083"
	}

	logger.Printf("Servicio de documentos iniciado en el puerto %s", port)

	//Iniciar el servidor
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal("Error al iniciar el servidor", err)
	}

}
