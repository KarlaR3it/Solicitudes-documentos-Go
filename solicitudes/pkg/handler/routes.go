package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kramirez/solicitudes/internal/solicitud"
)

func SetupRoutes(endpoints *solicitud.Endpoint) *gin.Engine {
	router := gin.Default()

	//Grupo de rutas para solicitudes
	solicitudGroup := router.Group("/solicitudes")
	{
		solicitudGroup.POST("", endpoints.Create)
		solicitudGroup.GET("", endpoints.GetAll)
		solicitudGroup.GET("/:id", endpoints.GetByID) // Obtiene solo la información básica
		solicitudGroup.GET("/:id/con-documentos", endpoints.GetByIDWithDocuments) // Obtiene la solicitud con sus documentos
		solicitudGroup.PATCH("/:id", endpoints.Update)
		solicitudGroup.DELETE("/:id", endpoints.Delete)
	}

	return router
}
