package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kramirez/documentos/internal/documento"
)

func SetupRoutes(endpoints *documento.Endpoint) *gin.Engine {
	router := gin.Default()

	//Grupo de rutas para documentos
	documentoGroup := router.Group("/documentos")
	{
		documentoGroup.POST("", endpoints.Create)
		documentoGroup.GET("", endpoints.GetAll)
		documentoGroup.GET("/:id", endpoints.GetByID)
		documentoGroup.PATCH("/:id", endpoints.Update)
		documentoGroup.DELETE("/:id", endpoints.Delete)
	}

	return router
}
