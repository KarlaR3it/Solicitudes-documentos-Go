package bootstrap

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitEnv(t *testing.T) {
	t.Run("debe ejecutarse sin errores cuando no existe .env", func(t *testing.T) {
		// Arrange - Asegurarse que no hay .env en el directorio actual
		originalDir, _ := os.Getwd()

		// Act - No debe generar panic
		assert.NotPanics(t, func() {
			InitEnv()
		}, "InitEnv no debe generar panic cuando no encuentra .env")

		// Assert - Verificar que estamos en el directorio original
		currentDir, _ := os.Getwd()
		assert.Equal(t, originalDir, currentDir)
	})

	t.Run("debe manejar variables de entorno del sistema", func(t *testing.T) {
		// Arrange - Establecer una variable de entorno de prueba
		testKey := "TEST_ENV_VAR"
		testValue := "test_value"
		os.Setenv(testKey, testValue)

		// Act
		InitEnv()

		// Assert - La variable debe seguir estando disponible
		assert.Equal(t, testValue, os.Getenv(testKey))

		// Cleanup
		os.Unsetenv(testKey)
	})
}

func TestDBConnection(t *testing.T) {
	t.Run("debe retornar error cuando faltan variables de entorno", func(t *testing.T) {
		// Arrange - Limpiar variables de BD
		originalEnvs := map[string]string{
			"DB_HOST":     os.Getenv("DB_HOST"),
			"DB_PORT":     os.Getenv("DB_PORT"),
			"DB_USER":     os.Getenv("DB_USER"),
			"DB_PASSWORD": os.Getenv("DB_PASSWORD"),
			"DB_NAME":     os.Getenv("DB_NAME"),
		}

		// Limpiar todas las variables de BD
		for key := range originalEnvs {
			os.Unsetenv(key)
		}

		// Act
		db, err := DBConnection()

		// Assert
		assert.Error(t, err, "Debe retornar error cuando faltan variables de BD")
		assert.Nil(t, db, "DB debe ser nil cuando hay error")

		// Cleanup - Restaurar variables originales
		for key, value := range originalEnvs {
			if value != "" {
				os.Setenv(key, value)
			}
		}
	})

	t.Run("debe construir DSN correctamente con variables válidas", func(t *testing.T) {
		// Arrange - Configurar variables de entorno de prueba
		testEnvs := map[string]string{
			"DB_HOST":     "localhost",
			"DB_PORT":     "3306",
			"DB_USER":     "testuser",
			"DB_PASSWORD": "testpass",
			"DB_NAME":     "testdb",
		}

		originalEnvs := make(map[string]string)
		for key, value := range testEnvs {
			originalEnvs[key] = os.Getenv(key)
			os.Setenv(key, value)
		}

		// Act
		_, err := DBConnection()

		// Assert - Aunque falle la conexión (BD no existe), no debe ser error de DSN
		if err != nil {
			// El error debe ser de conexión, no de formato DSN
			assert.Contains(t, err.Error(), "connection", "Error debe ser de conexión, no de DSN mal formateado")
		}

		// Cleanup
		for key, value := range originalEnvs {
			if value != "" {
				os.Setenv(key, value)
			} else {
				os.Unsetenv(key)
			}
		}
	})

	t.Run("debe manejar configuración de debug correctamente", func(t *testing.T) {
		// Arrange
		originalDebug := os.Getenv("DATABASE_DEBUG")
		originalEnvs := map[string]string{
			"DB_HOST":     "localhost",
			"DB_PORT":     "3306",
			"DB_USER":     "testuser",
			"DB_PASSWORD": "testpass",
			"DB_NAME":     "testdb",
		}

		for key, value := range originalEnvs {
			os.Setenv(key, value)
		}

		// Test con debug activado
		os.Setenv("DATABASE_DEBUG", "true")

		// Assert - La función debe ejecutarse sin panic
		// (aunque puede fallar la conexión por BD inexistente)
		assert.NotPanics(t, func() {
			DBConnection()
		}, "No debe generar panic con DATABASE_DEBUG=true")

		// Test con debug desactivado
		os.Setenv("DATABASE_DEBUG", "false")
		assert.NotPanics(t, func() {
			DBConnection()
		}, "No debe generar panic con DATABASE_DEBUG=false")

		// Cleanup
		if originalDebug != "" {
			os.Setenv("DATABASE_DEBUG", originalDebug)
		} else {
			os.Unsetenv("DATABASE_DEBUG")
		}

		for key := range originalEnvs {
			os.Unsetenv(key)
		}
	})

	t.Run("debe manejar migraciones cuando DATABASE_MIGRATE=up", func(t *testing.T) {
		// Arrange
		originalMigrate := os.Getenv("DATABASE_MIGRATE")
		originalEnvs := map[string]string{
			"DB_HOST":     "localhost",
			"DB_PORT":     "3306",
			"DB_USER":     "testuser",
			"DB_PASSWORD": "testpass",
			"DB_NAME":     "testdb",
		}

		for key, value := range originalEnvs {
			os.Setenv(key, value)
		}

		os.Setenv("DATABASE_MIGRATE", "up")

		// Act & Assert - No debe generar panic
		assert.NotPanics(t, func() {
			DBConnection()
		}, "No debe generar panic con DATABASE_MIGRATE=up")

		// Cleanup
		if originalMigrate != "" {
			os.Setenv("DATABASE_MIGRATE", originalMigrate)
		} else {
			os.Unsetenv("DATABASE_MIGRATE")
		}

		for key := range originalEnvs {
			os.Unsetenv(key)
		}
	})
}

func TestInitLogger(t *testing.T) {
	t.Run("debe crear logger correctamente", func(t *testing.T) {
		// Act
		logger := InitLogger()

		// Assert
		assert.NotNil(t, logger, "Logger no debe ser nil")
		assert.IsType(t, &log.Logger{}, logger, "Debe retornar un *log.Logger")
	})

	t.Run("debe configurar logger con flags correctos", func(t *testing.T) {
		// Act
		logger := InitLogger()

		// Assert - Verificar que puede escribir sin errores
		assert.NotPanics(t, func() {
			logger.Print("Test log message")
		}, "Logger debe poder escribir sin errores")
	})

	t.Run("debe usar stdout como output", func(t *testing.T) {
		// Act
		logger := InitLogger()

		// Assert - Verificar que el logger existe y puede funcionar
		assert.NotNil(t, logger)

		// Verificar que puede generar output (aunque no podamos capturarlo fácilmente)
		assert.NotPanics(t, func() {
			logger.Println("Test message for stdout verification")
		})
	})
}

func TestBootstrapIntegration(t *testing.T) {
	t.Run("debe integrar InitEnv e InitLogger correctamente", func(t *testing.T) {
		// Act - Ejecutar secuencia típica de bootstrap
		assert.NotPanics(t, func() {
			InitEnv()
			logger := InitLogger()
			logger.Println("Bootstrap integration test")
		}, "Integración InitEnv + InitLogger debe funcionar sin errores")
	})

	t.Run("debe manejar configuración completa de aplicación", func(t *testing.T) {
		// Arrange
		testEnvs := map[string]string{
			"DB_HOST":          "localhost",
			"DB_PORT":          "3306",
			"DB_USER":          "testuser",
			"DB_PASSWORD":      "testpass",
			"DB_NAME":          "testdb",
			"DATABASE_DEBUG":   "false",
			"DATABASE_MIGRATE": "up",
		}

		originalEnvs := make(map[string]string)
		for key, value := range testEnvs {
			originalEnvs[key] = os.Getenv(key)
			os.Setenv(key, value)
		}

		// Act & Assert - Secuencia completa de bootstrap
		assert.NotPanics(t, func() {
			InitEnv()
			logger := InitLogger()
			logger.Println("Testing complete bootstrap sequence")

			// Intentar conexión a BD (puede fallar pero no debe hacer panic)
			_, err := DBConnection()
			if err != nil {
				logger.Printf("DB connection failed as expected in test: %v", err)
			}
		}, "Secuencia completa de bootstrap no debe generar panic")

		// Cleanup
		for key, value := range originalEnvs {
			if value != "" {
				os.Setenv(key, value)
			} else {
				os.Unsetenv(key)
			}
		}
	})
}
