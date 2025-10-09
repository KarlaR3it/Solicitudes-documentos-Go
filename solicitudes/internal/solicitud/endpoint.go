package solicitud

import (
	"encoding/json"
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

// Create maneja POST /solicitudes
func (e *Endpoint) Create(c *gin.Context) {
	var req CreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	solicitud, err := e.service.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, solicitud)
}

// GetAll maneja GET /solicitudes
func (e *Endpoint) GetAll(c *gin.Context) {
	filters := GetAllReq{
		Titulo:           c.Query("titulo"),
		Estado:           c.Query("estado"),
		Area:             c.Query("area"),
		Pais:             c.Query("pais"),
		ModalidadTrabajo: c.Query("modalidadTrabajo"),
		TipoServicio:     c.Query("tipoServicio"),
	}

	// Convertir RentaDesde
	if rentaDesde := c.Query("rentaDesde"); rentaDesde != "" {
		if rd, err := strconv.Atoi(rentaDesde); err == nil {
			filters.RentaDesde = rd
		}
	}

	// Convertir RentaHasta
	if rentaHasta := c.Query("rentaHasta"); rentaHasta != "" {
		if rh, err := strconv.Atoi(rentaHasta); err == nil {
			filters.RentaHasta = rh
		}
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

	solicitudes, err := e.service.GetAll(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, solicitudes)
}

// GetByID maneja GET /solicitudes/:id
func (e *Endpoint) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	solicitud, err := e.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Solicitud no encontrada"})
		return
	}

	c.JSON(http.StatusOK, solicitud)
}

// GetByIDWithDocuments maneja GET /solicitudes/:id/con-documentos
func (e *Endpoint) GetByIDWithDocuments(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	solicitud, err := e.service.GetByIDWithDocuments(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Solicitud no encontrada"})
		return
	}

	c.JSON(http.StatusOK, solicitud)
}

// Update maneja PATCH /solicitudes/:id
func (e *Endpoint) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Leer el body una sola vez
	bodyBytes, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar campos no permitidos
	var rawBody map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &rawBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validar que no se intenten actualizar campos no permitidos
	forbiddenFields := []string{"usuario_id", "id", "created_at", "updated_at"}
	for _, field := range forbiddenFields {
		if _, exists := rawBody[field]; exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Campo '" + field + "' no puede ser actualizado"})
			return
		}
	}

	// Parsear con la estructura correcta
	var req UpdateReq
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := e.service.Update(c.Request.Context(), uint(id), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Solicitud actualizada exitosamente"})
}

// Delete, maneja DELETE /solicitudes/:id
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
	c.JSON(http.StatusOK, gin.H{"message": "Solicitud eliminada exitosamente"})
}
