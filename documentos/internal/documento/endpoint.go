package documento

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Endpoint struct {
	service Service
}

func NewEndpoint(service Service) *Endpoint {
	return &Endpoint{service: service}
}

// Create maneja POST /documentos
func (e *Endpoint) Create(c *gin.Context) {
	var req CreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	documento, err := e.service.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, documento)
}

// GetAll maneja GET /documentos
func (e *Endpoint) GetAll(c *gin.Context) {
	filters := GetAllReq{
		Extension:     c.Query("extension"),
		NombreArchivo: c.Query("nombre_archivo"),
	}

	//Paginacion
	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			filters.Limit = l
		}
	}
	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			filters.Page = p
		}
	}

	documentos, err := e.service.GetAll(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, documentos)
}

// GetByID maneja GET /documentos/:id
func (e *Endpoint) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	documento, err := e.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Documento no encontrado"})
		return
	}
	c.JSON(http.StatusOK, documento)
}

// Update maneja PATCH /documentos/:id
func (e *Endpoint) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req UpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := e.service.Update(c.Request.Context(), uint(id), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Documento actualizado exitosamente"})
}

// Delete maneja DELETE /documentos/:id
func (e *Endpoint) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := e.service.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Documento eliminado exitosamente"})
}
