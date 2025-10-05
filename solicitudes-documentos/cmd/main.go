package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kramirez/solicitudes-documentos/internal/solicitud"
	"github.com/kramirez/solicitudes-documentos/pkg/bootstrap"
	"github.com/kramirez/solicitudes-documentos/pkg/handler"
	"github.com/kramirez/solicitudes-documentos/pkg/httpclient"
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

	//Inicializar cliente HTTP para el servicio de usuarios
	usuariosServiceURL := os.Getenv("USUARIOS_SERVICE_URL")
	if usuariosServiceURL == "" {
		usuariosServiceURL = "http://localhost:8081"
	}
	usuarioClient := httpclient.NewUsuarioClient(usuariosServiceURL)
	logger.Printf("Cliente de usuarios configurado: %s", usuariosServiceURL)

	//Inicializar capas
	repo := solicitud.NewRepository(db)
	service := solicitud.NewService(repo, logger, usuarioClient)
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
