package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kramirez/solicitudes/internal/solicitud"
)

func SetupRoutes(endpoints *solicitud.Endpoint) *gin.Engine {
	router := gin.Default()

	// Configurar CORS (permitir todos los orígenes - solo para desarrollo)
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true, //Esto es solo para ambiente de desarrollo, para producción se debe configurar los orígenes
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	//Grupo de rutas para solicitudes
	solicitudGroup := router.Group("/solicitudes")
	{
		solicitudGroup.POST("", endpoints.Create)
		solicitudGroup.GET("", endpoints.GetAll)
		solicitudGroup.GET("/:id", endpoints.GetByID)                             // Obtiene solo la información básica
		solicitudGroup.GET("/:id/con-documentos", endpoints.GetByIDWithDocuments) // Obtiene la solicitud con sus documentos
		solicitudGroup.PATCH("/:id", endpoints.Update)
		solicitudGroup.DELETE("/:id", endpoints.Delete)
	}

	return router
}
