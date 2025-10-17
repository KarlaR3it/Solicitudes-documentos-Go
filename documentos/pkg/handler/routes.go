package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kramirez/documentos/internal/documento"
)

func SetupRoutes(endpoints *documento.Endpoint) *gin.Engine {
	router := gin.Default()

	// Configurar CORS (permitir todos los orígenes - solo para desarrollo)
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true, //Esto es solo para ambiente de desarrollo, para producción se debe configurar los orígenes permitidos
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	//Grupo de rutas para documentos
	documentoGroup := router.Group("/documentos")
	{
		documentoGroup.POST("", endpoints.Create)
		documentoGroup.GET("", endpoints.GetAll)
		documentoGroup.GET("/:id", endpoints.GetByID)
		documentoGroup.PATCH("/:id", endpoints.Update)
		documentoGroup.DELETE("/:id", endpoints.Delete)
		documentoGroup.DELETE("/solicitud/:solicitud_id", endpoints.DeleteBySolicitudID)
	}

	return router
}
