package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kramirez/solicitudes/internal/solicitud"
	"github.com/stretchr/testify/assert"
)

func TestSetupRoutes(t *testing.T) {
	// Configurar modo de prueba para Gin
	gin.SetMode(gin.TestMode)

	// Crear un mock endpoint para las pruebas
	mockEndpoint := &solicitud.Endpoint{}

	t.Run("debe configurar todas las rutas correctamente", func(t *testing.T) {
		// Arrange
		router := SetupRoutes(mockEndpoint)

		// Verificar que el router se haya creado
		assert.NotNil(t, router)

		// Verificar estructura de rutas usando el router info
		routes := router.Routes()

		// Mapear rutas esperadas (usando slice porque hay múltiples GET)
		expectedRoutes := []struct {
			method string
			path   string
		}{
			{"POST", "/solicitudes"},
			{"GET", "/solicitudes"},
			{"GET", "/solicitudes/:id"},
			{"GET", "/solicitudes/:id/con-documentos"},
			{"PATCH", "/solicitudes/:id"},
			{"DELETE", "/solicitudes/:id"},
		}

		// Verificar que se registraron las rutas correctas
		assert.GreaterOrEqual(t, len(routes), len(expectedRoutes))

		// Verificar que todas las rutas esperadas están presentes
		routeMap := make(map[string]bool)
		for _, route := range routes {
			key := route.Method + ":" + route.Path
			routeMap[key] = true
		}

		// Validar rutas específicas
		for _, expected := range expectedRoutes {
			key := expected.method + ":" + expected.path
			assert.True(t, routeMap[key], "Ruta %s %s debe existir", expected.method, expected.path)
		}
	})

	t.Run("debe responder a rutas POST correctamente", func(t *testing.T) {
		// Arrange
		router := SetupRoutes(mockEndpoint)

		req := httptest.NewRequest("POST", "/solicitudes", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert - Verificar que la ruta existe (no 404)
		assert.NotEqual(t, http.StatusNotFound, w.Code, "La ruta POST /solicitudes debe existir")
	})

	t.Run("debe responder a rutas GET correctamente", func(t *testing.T) {
		// Arrange
		router := SetupRoutes(mockEndpoint)

		// Test GET /solicitudes
		req := httptest.NewRequest("GET", "/solicitudes", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert - Verificar que la ruta existe (no 404)
		assert.NotEqual(t, http.StatusNotFound, w.Code, "La ruta GET /solicitudes debe existir")
	})

	t.Run("debe responder a rutas GET con ID correctamente", func(t *testing.T) {
		// Arrange
		router := SetupRoutes(mockEndpoint)

		// Test GET /solicitudes/1
		req := httptest.NewRequest("GET", "/solicitudes/1", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert - Verificar que la ruta existe (no 404)
		assert.NotEqual(t, http.StatusNotFound, w.Code, "La ruta GET /solicitudes/:id debe existir")
	})

	t.Run("debe responder a rutas con documentos correctamente", func(t *testing.T) {
		// Arrange
		router := SetupRoutes(mockEndpoint)

		// Test GET /solicitudes/1/con-documentos
		req := httptest.NewRequest("GET", "/solicitudes/1/con-documentos", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert - Verificar que la ruta existe (no 404)
		assert.NotEqual(t, http.StatusNotFound, w.Code, "La ruta GET /solicitudes/:id/con-documentos debe existir")
	})

	t.Run("debe responder a rutas PATCH correctamente", func(t *testing.T) {
		// Arrange
		router := SetupRoutes(mockEndpoint)

		req := httptest.NewRequest("PATCH", "/solicitudes/1", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert - Verificar que la ruta existe (no 404)
		assert.NotEqual(t, http.StatusNotFound, w.Code, "La ruta PATCH /solicitudes/:id debe existir")
	})

	t.Run("debe responder a rutas DELETE correctamente", func(t *testing.T) {
		// Arrange
		router := SetupRoutes(mockEndpoint)

		req := httptest.NewRequest("DELETE", "/solicitudes/1", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert - Verificar que la ruta existe (no 404)
		assert.NotEqual(t, http.StatusNotFound, w.Code, "La ruta DELETE /solicitudes/:id debe existir")
	})

	t.Run("debe retornar 404 para rutas inexistentes", func(t *testing.T) {
		// Arrange
		router := SetupRoutes(mockEndpoint)

		req := httptest.NewRequest("GET", "/ruta-inexistente", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code, "Rutas inexistentes deben retornar 404")
	})

	t.Run("debe agrupar correctamente las rutas bajo /solicitudes", func(t *testing.T) {
		// Arrange
		router := SetupRoutes(mockEndpoint)
		routes := router.Routes()

		// Act & Assert - Verificar que todas las rutas están bajo el grupo /solicitudes
		solicitudRoutes := 0
		for _, route := range routes {
			if strings.HasPrefix(route.Path, "/solicitudes") {
				solicitudRoutes++
			}
		}

		// Debe haber al menos 6 rutas bajo /solicitudes
		assert.GreaterOrEqual(t, solicitudRoutes, 6, "Debe haber al menos 6 rutas bajo /solicitudes")
	})
}
