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

func TestNewUsuarioClient(t *testing.T) {
	t.Run("debe crear cliente correctamente con URL válida", func(t *testing.T) {
		// Arrange
		baseURL := "http://localhost:8081"

		// Act
		client := NewUsuarioClient(baseURL)

		// Assert
		assert.NotNil(t, client)
		assert.Equal(t, baseURL, client.baseURL)
		assert.NotNil(t, client.httpClient)
		assert.Equal(t, 10*time.Second, client.httpClient.Timeout)
	})

	t.Run("debe crear cliente con URL vacía", func(t *testing.T) {
		// Arrange
		baseURL := ""

		// Act
		client := NewUsuarioClient(baseURL)

		// Assert
		assert.NotNil(t, client)
		assert.Equal(t, "", client.baseURL)
		assert.NotNil(t, client.httpClient)
	})

	t.Run("debe configurar timeout correctamente", func(t *testing.T) {
		// Arrange & Act
		client := NewUsuarioClient("http://test.com")

		// Assert
		assert.Equal(t, 10*time.Second, client.httpClient.Timeout)
	})
}

func TestUsuarioClient_ValidarUsuario(t *testing.T) {
	t.Run("debe retornar true cuando usuario existe", func(t *testing.T) {
		// Arrange - Mock server que retorna 200 OK
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verificar método y URL
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/usuarios/123")

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(UsuarioResponse{ID: 123, NombreUsuario: "testuser"})
		}))
		defer server.Close()

		client := NewUsuarioClient(server.URL)

		// Act
		existe, err := client.ValidarUsuario(123)

		// Assert
		assert.NoError(t, err)
		assert.True(t, existe)
	})

	t.Run("debe retornar false cuando usuario no existe", func(t *testing.T) {
		// Arrange - Mock server que retorna 404
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/usuarios/999")

			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		client := NewUsuarioClient(server.URL)

		// Act
		existe, err := client.ValidarUsuario(999)

		// Assert
		assert.NoError(t, err)
		assert.False(t, existe)
	})

	t.Run("debe manejar error de conexión", func(t *testing.T) {
		// Arrange - URL inválida
		client := NewUsuarioClient("http://servidor-inexistente:9999")

		// Act
		existe, err := client.ValidarUsuario(123)

		// Assert
		assert.Error(t, err)
		assert.False(t, existe)
		assert.Contains(t, err.Error(), "error al conectar con el servicio de usuarios")
	})

	t.Run("debe manejar códigos de estado de error distintos a 404", func(t *testing.T) {
		testCases := []struct {
			statusCode   int
			expectError  bool
			expectExists bool
		}{
			{http.StatusOK, false, true},
			{http.StatusNotFound, false, false},
			{http.StatusBadRequest, true, false},
			{http.StatusInternalServerError, true, false},
			{http.StatusForbidden, true, false},
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("status_%d", tc.statusCode), func(t *testing.T) {
				// Arrange
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tc.statusCode)
					if tc.statusCode == http.StatusOK {
						json.NewEncoder(w).Encode(UsuarioResponse{ID: 123})
					}
				}))
				defer server.Close()

				client := NewUsuarioClient(server.URL)

				// Act
				existe, err := client.ValidarUsuario(123)

				// Assert
				if tc.expectError {
					assert.Error(t, err)
					assert.Contains(t, err.Error(), "error del servicio de usuarios")
				} else {
					assert.NoError(t, err)
				}
				assert.Equal(t, tc.expectExists, existe)
			})
		}
	})

	t.Run("debe construir URL correctamente", func(t *testing.T) {
		// Arrange - Mock server para capturar URL
		var capturedURL string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedURL = r.URL.Path
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client := NewUsuarioClient(server.URL)

		// Act
		client.ValidarUsuario(42)

		// Assert
		assert.Equal(t, "/usuarios/42", capturedURL)
	})
}

func TestUsuarioClient_ObtenerUsuario(t *testing.T) {
	t.Run("debe obtener usuario exitosamente", func(t *testing.T) {
		// Arrange - Mock usuario
		mockUsuario := UsuarioResponse{
			ID:            123,
			NombreUsuario: "johndoe",
			EmailUsuario:  "john@example.com",
			Nombre:        "John",
			Apellidos:     "Doe",
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verificar método y URL
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/usuarios/123")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(mockUsuario)
		}))
		defer server.Close()

		client := NewUsuarioClient(server.URL)

		// Act
		usuario, err := client.ObtenerUsuario(123)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, usuario)
		assert.Equal(t, uint(123), usuario.ID)
		assert.Equal(t, "johndoe", usuario.NombreUsuario)
		assert.Equal(t, "john@example.com", usuario.EmailUsuario)
		assert.Equal(t, "John", usuario.Nombre)
		assert.Equal(t, "Doe", usuario.Apellidos)
	})

	t.Run("debe retornar error cuando usuario no existe", func(t *testing.T) {
		// Arrange - Mock server que retorna 404
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/usuarios/999")

			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		client := NewUsuarioClient(server.URL)

		// Act
		usuario, err := client.ObtenerUsuario(999)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, usuario)
		assert.Contains(t, err.Error(), "usuario no encontrado")
	})

	t.Run("debe manejar error de conexión", func(t *testing.T) {
		// Arrange - URL inválida
		client := NewUsuarioClient("http://servidor-inexistente:9999")

		// Act
		usuario, err := client.ObtenerUsuario(123)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, usuario)
		assert.Contains(t, err.Error(), "error al conectar con el servicio de usuarios")
	})

	t.Run("debe manejar JSON inválido en respuesta", func(t *testing.T) {
		// Arrange - Mock server con JSON malformado
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{json malformado`))
		}))
		defer server.Close()

		client := NewUsuarioClient(server.URL)

		// Act
		usuario, err := client.ObtenerUsuario(123)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, usuario)
		assert.Contains(t, err.Error(), "error al decodificar respuesta")
	})

	t.Run("debe manejar diferentes códigos de estado HTTP", func(t *testing.T) {
		testCases := []struct {
			statusCode    int
			expectError   bool
			errorContains string
		}{
			{http.StatusOK, false, ""},
			{http.StatusNotFound, true, "usuario no encontrado"},
			{http.StatusBadRequest, true, "error del servicio de usuarios: status 400"},
			{http.StatusInternalServerError, true, "error del servicio de usuarios: status 500"},
			{http.StatusForbidden, true, "error del servicio de usuarios: status 403"},
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("status_%d", tc.statusCode), func(t *testing.T) {
				// Arrange
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tc.statusCode)
					if tc.statusCode == http.StatusOK {
						json.NewEncoder(w).Encode(UsuarioResponse{ID: 123})
					}
				}))
				defer server.Close()

				client := NewUsuarioClient(server.URL)

				// Act
				usuario, err := client.ObtenerUsuario(123)

				// Assert
				if tc.expectError {
					assert.Error(t, err)
					assert.Nil(t, usuario)
					assert.Contains(t, err.Error(), tc.errorContains)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, usuario)
				}
			})
		}
	})

	t.Run("debe construir URL correctamente", func(t *testing.T) {
		// Arrange - Mock server para capturar URL
		var capturedURL string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedURL = r.URL.Path
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(UsuarioResponse{ID: 42})
		}))
		defer server.Close()

		client := NewUsuarioClient(server.URL)

		// Act
		client.ObtenerUsuario(42)

		// Assert
		assert.Equal(t, "/usuarios/42", capturedURL)
	})

	t.Run("debe decodificar respuesta JSON correctamente", func(t *testing.T) {
		// Arrange - Usuario con todos los campos
		mockUsuario := UsuarioResponse{
			ID:            456,
			NombreUsuario: "jane_doe",
			EmailUsuario:  "jane.doe@company.com",
			Nombre:        "Jane Elizabeth",
			Apellidos:     "Doe Smith",
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(mockUsuario)
		}))
		defer server.Close()

		client := NewUsuarioClient(server.URL)

		// Act
		usuario, err := client.ObtenerUsuario(456)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, usuario)
		assert.Equal(t, mockUsuario.ID, usuario.ID)
		assert.Equal(t, mockUsuario.NombreUsuario, usuario.NombreUsuario)
		assert.Equal(t, mockUsuario.EmailUsuario, usuario.EmailUsuario)
		assert.Equal(t, mockUsuario.Nombre, usuario.Nombre)
		assert.Equal(t, mockUsuario.Apellidos, usuario.Apellidos)
	})
}

func TestUsuarioResponse(t *testing.T) {
	t.Run("debe crear estructura UsuarioResponse correctamente", func(t *testing.T) {
		// Arrange & Act
		usuario := UsuarioResponse{
			ID:            1,
			NombreUsuario: "test",
			EmailUsuario:  "test@example.com",
			Nombre:        "Test",
			Apellidos:     "User",
		}

		// Assert
		assert.Equal(t, uint(1), usuario.ID)
		assert.Equal(t, "test", usuario.NombreUsuario)
		assert.Equal(t, "test@example.com", usuario.EmailUsuario)
		assert.Equal(t, "Test", usuario.Nombre)
		assert.Equal(t, "User", usuario.Apellidos)
	})

	t.Run("debe serializar/deserializar JSON correctamente", func(t *testing.T) {
		// Arrange
		original := UsuarioResponse{
			ID:            789,
			NombreUsuario: "serialization_test",
			EmailUsuario:  "serial@test.com",
			Nombre:        "Serial",
			Apellidos:     "Test",
		}

		// Act - Serializar
		jsonData, err := json.Marshal(original)
		assert.NoError(t, err)

		// Act - Deserializar
		var deserialized UsuarioResponse
		err = json.Unmarshal(jsonData, &deserialized)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, original.ID, deserialized.ID)
		assert.Equal(t, original.NombreUsuario, deserialized.NombreUsuario)
		assert.Equal(t, original.EmailUsuario, deserialized.EmailUsuario)
		assert.Equal(t, original.Nombre, deserialized.Nombre)
		assert.Equal(t, original.Apellidos, deserialized.Apellidos)
	})
}
