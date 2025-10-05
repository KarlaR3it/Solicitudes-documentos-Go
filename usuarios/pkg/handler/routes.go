package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kramirez/usuarios-service/internal/usuario"
)

func SetupRoutes(endpoint *usuario.Endpoint) *gin.Engine {
	router := gin.Default()

	// Grupo de rutas para usuarios
	usuariosGroup := router.Group("/usuarios")
	{
		usuariosGroup.POST("", endpoint.Create)
		usuariosGroup.GET("", endpoint.GetAll)
		usuariosGroup.GET("/:id", endpoint.GetByID)
		usuariosGroup.PATCH("/:id", endpoint.Update)
		usuariosGroup.DELETE("/:id", endpoint.Delete)
	}

	return router
}
