package documento

import (
	"encoding/json"
	"fmt"
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
	// Configurar decoder para rechazar campos desconocidos
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()
	
	if err := decoder.Decode(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Validar campos requeridos manualmente
	if req.Extension == "" || req.NombreArchivo == "" || req.SolicitudID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Todos los campos son requeridos"})
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

	// Convertir solicitud_id
	if solicitudID := c.Query("solicitud_id"); solicitudID != "" {
		if sid, err := strconv.Atoi(solicitudID); err == nil {
			filters.SolicitudID = uint(sid)
		}
	}

	// Paginacion
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

	// Validar campos prohibidos y desconocidos
	var rawReq map[string]interface{}
	if err := c.ShouldBindJSON(&rawReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Lista de campos que no se pueden actualizar
	forbiddenFields := []string{"id", "created_at", "updated_at", "solicitud_id"}
	for _, field := range forbiddenFields {
		if _, exists := rawReq[field]; exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Campo '%s' no puede ser actualizado", field)})
			return
		}
	}

	// Lista de campos permitidos
	allowedFields := map[string]bool{
		"extension":      true,
		"nombre_archivo": true,
	}

	// Validar que no haya campos desconocidos
	for field := range rawReq {
		if !allowedFields[field] {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Campo '%s' no es válido", field)})
			return
		}
	}

	// Convertir a UpdateReq
	var req UpdateReq
	if rawReq["extension"] != nil {
		if ext, ok := rawReq["extension"].(string); ok {
			req.Extension = &ext
		}
	}
	if rawReq["nombre_archivo"] != nil {
		if nombre, ok := rawReq["nombre_archivo"].(string); ok {
			req.NombreArchivo = &nombre
		}
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

// DeleteBySolicitudID maneja DELETE /documentos/solicitud/:solicitud_id
func (e *Endpoint) DeleteBySolicitudID(c *gin.Context) {
	solicitudID, err := strconv.ParseUint(c.Param("solicitud_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de solicitud inválido"})
		return
	}

	if err := e.service.DeleteBySolicitudID(c.Request.Context(), uint(solicitudID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Documentos eliminados exitosamente"})
}
