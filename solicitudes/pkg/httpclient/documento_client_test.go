package httpclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDocumentoDTO_toSolicitudDocumento(t *testing.T) {
	t.Run("debe convertir DocumentoDTO a Documento correctamente", func(t *testing.T) {
		// Arrange
		dto := DocumentoDTO{
			ID:            1,
			Extension:     "pdf",
			NombreArchivo: "documento_test.pdf",
			SolicitudID:   100,
		}

		// Act
		documento := dto.toSolicitudDocumento()

		// Assert
		assert.Equal(t, uint(1), documento.ID)
		assert.Equal(t, "pdf", documento.Extension)
		assert.Equal(t, "documento_test.pdf", documento.NombreArchivo)
	})

	t.Run("debe manejar valores vacíos correctamente", func(t *testing.T) {
		// Arrange
		dto := DocumentoDTO{}

		// Act
		documento := dto.toSolicitudDocumento()

		// Assert
		assert.Equal(t, uint(0), documento.ID)
		assert.Equal(t, "", documento.Extension)
		assert.Equal(t, "", documento.NombreArchivo)
	})

	t.Run("debe preservar todos los campos durante la conversión", func(t *testing.T) {
		// Arrange
		dto := DocumentoDTO{
			ID:            999,
			Extension:     "jpg",
			NombreArchivo: "imagen_compleja_nombre.jpg",
			SolicitudID:   555,
		}

		// Act
		documento := dto.toSolicitudDocumento()

		// Assert
		assert.Equal(t, dto.ID, documento.ID)
		assert.Equal(t, dto.Extension, documento.Extension)
		assert.Equal(t, dto.NombreArchivo, documento.NombreArchivo)
	})
}

func TestNewDocumentoClient(t *testing.T) {
	t.Run("debe crear cliente correctamente con URL válida", func(t *testing.T) {
		// Arrange
		baseURL := "http://localhost:8080"

		// Act
		client := NewDocumentoClient(baseURL)

		// Assert
		assert.NotNil(t, client)
		assert.Equal(t, baseURL, client.baseURL)
		assert.NotNil(t, client.client)
		assert.Equal(t, 10*time.Second, client.client.Timeout)
	})

	t.Run("debe crear cliente con URL vacía", func(t *testing.T) {
		// Arrange
		baseURL := ""

		// Act
		client := NewDocumentoClient(baseURL)

		// Assert
		assert.NotNil(t, client)
		assert.Equal(t, "", client.baseURL)
		assert.NotNil(t, client.client)
	})

	t.Run("debe configurar timeout correctamente", func(t *testing.T) {
		// Arrange & Act
		client := NewDocumentoClient("http://test.com")

		// Assert
		assert.Equal(t, 10*time.Second, client.client.Timeout)
	})
}

func TestDocumentoClient_GetBySolicitudID(t *testing.T) {
	t.Run("debe obtener documentos exitosamente", func(t *testing.T) {
		// Arrange - Mock server
		mockDocumentos := []DocumentoDTO{
			{ID: 1, Extension: "pdf", NombreArchivo: "doc1.pdf", SolicitudID: 100},
			{ID: 2, Extension: "jpg", NombreArchivo: "img1.jpg", SolicitudID: 100},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verificar método y URL
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/documentos")
			assert.Equal(t, "100", r.URL.Query().Get("solicitud_id"))
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			// Responder con documentos mock
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(mockDocumentos)
		}))
		defer server.Close()

		client := NewDocumentoClient(server.URL)

		// Act
		documentos, err := client.GetBySolicitudID(100)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, documentos, 2)
		assert.Equal(t, uint(1), documentos[0].ID)
		assert.Equal(t, "pdf", documentos[0].Extension)
		assert.Equal(t, "doc1.pdf", documentos[0].NombreArchivo)
		assert.Equal(t, uint(2), documentos[1].ID)
		assert.Equal(t, "jpg", documentos[1].Extension)
	})

	t.Run("debe manejar respuesta vacía correctamente", func(t *testing.T) {
		// Arrange - Mock server con respuesta vacía
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]DocumentoDTO{})
		}))
		defer server.Close()

		client := NewDocumentoClient(server.URL)

		// Act
		documentos, err := client.GetBySolicitudID(999)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, documentos, 0)
	})

	t.Run("debe manejar error de conexión", func(t *testing.T) {
		// Arrange - URL inválida
		client := NewDocumentoClient("http://servidor-inexistente:9999")

		// Act
		documentos, err := client.GetBySolicitudID(100)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, documentos)
		assert.Contains(t, err.Error(), "error al conectar con el servicio")
	})

	t.Run("debe manejar código de estado HTTP de error", func(t *testing.T) {
		// Arrange - Mock server que retorna 404
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		client := NewDocumentoClient(server.URL)

		// Act
		documentos, err := client.GetBySolicitudID(100)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, documentos)
		assert.Contains(t, err.Error(), "error al obtener documentos: status 404")
	})

	t.Run("debe manejar JSON inválido en respuesta", func(t *testing.T) {
		// Arrange - Mock server con JSON malformado
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{json malformado`))
		}))
		defer server.Close()

		client := NewDocumentoClient(server.URL)

		// Act
		documentos, err := client.GetBySolicitudID(100)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, documentos)
		assert.Contains(t, err.Error(), "error al decodificar respuesta")
	})

	t.Run("debe construir URL correctamente", func(t *testing.T) {
		// Arrange - Mock server para verificar URL
		var capturedURL string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedURL = r.URL.String()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]DocumentoDTO{})
		}))
		defer server.Close()

		client := NewDocumentoClient(server.URL)

		// Act
		client.GetBySolicitudID(42)

		// Assert
		expectedPath := "/documentos?solicitud_id=42"
		assert.Equal(t, expectedPath, capturedURL)
	})

	t.Run("debe manejar diferentes códigos de estado HTTP", func(t *testing.T) {
		testCases := []struct {
			statusCode   int
			expectError  bool
			errorMessage string
		}{
			{http.StatusOK, false, ""},
			{http.StatusBadRequest, true, "status 400"},
			{http.StatusInternalServerError, true, "status 500"},
			{http.StatusForbidden, true, "status 403"},
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("status_%d", tc.statusCode), func(t *testing.T) {
				// Arrange
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tc.statusCode)
					if tc.statusCode == http.StatusOK {
						json.NewEncoder(w).Encode([]DocumentoDTO{})
					}
				}))
				defer server.Close()

				client := NewDocumentoClient(server.URL)

				// Act
				documentos, err := client.GetBySolicitudID(100)

				// Assert
				if tc.expectError {
					assert.Error(t, err)
					assert.Nil(t, documentos)
					assert.Contains(t, err.Error(), tc.errorMessage)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, documentos)
				}
			})
		}
	})

	t.Run("debe configurar headers correctamente", func(t *testing.T) {
		// Arrange - Mock server para verificar headers
		var capturedHeaders http.Header
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedHeaders = r.Header.Clone()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]DocumentoDTO{})
		}))
		defer server.Close()

		client := NewDocumentoClient(server.URL)

		// Act
		client.GetBySolicitudID(100)

		// Assert
		assert.Equal(t, "application/json", capturedHeaders.Get("Content-Type"))
	})
}
