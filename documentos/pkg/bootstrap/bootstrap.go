package bootstrap

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kramirez/documentos/internal/documento"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitEnv() {
	// Mostrar el directorio de trabajo actual
	if wd, err := os.Getwd(); err == nil {
		log.Printf("Directorio de trabajo actual: %s\n", wd)
	}

	// Intentar cargar .env desde varias ubicaciones posibles
	paths := []string{
		".env",                    // Directorio actual
		"../documentos/.env",      // Ruta relativa desde el raíz del módulo
		"C:/PruebasEureka/solicitudes-Go/documentos/.env",  // Ruta absoluta con barras normales
	}

	var loaded bool
	for _, path := range paths {
		log.Printf("Intentando cargar: %s\n", path)
		err := godotenv.Load(path)
		if err == nil {
			log.Printf("Archivo .env cargado exitosamente desde: %s\n", path)
			loaded = true
			break
		} else {
			log.Printf("No se pudo cargar el archivo %s: %v\n", path, err)
		}
	}

	if !loaded {
		log.Println("No se encontró archivo .env, usando variables de entorno del sistema")
	}

	// Depuración: Mostrar variables de entorno relevantes
	envVars := []string{"DB_HOST", "DB_PORT", "DB_NAME", "DB_USER", "DB_PASSWORD"}
	log.Println("Valores de las variables de entorno:")
	for _, envVar := range envVars {
		log.Printf("%s=%s\n", envVar, os.Getenv(envVar))
	}
}

func DBConnection() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))
	
	log.Printf("Intentando conectar a la base de datos con DSN: %s:%s@tcp(%s:%s)/%s\n", 
		os.Getenv("DB_USER"),
		"****", // No mostramos la contraseña por seguridad
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Error al conectar a la base de datos: %v\n", err)
		return nil, err
	}

	if os.Getenv("DATABASE_DEBUG") == "true" {
		db = db.Debug()
	}

	// Auto-migrate si DATABASE_MIGRATE está en "up"
	if os.Getenv("DATABASE_MIGRATE") == "up" {
		if err := db.AutoMigrate(&documento.Documento{}); err != nil {
			return nil, fmt.Errorf("error al realizar migraciones: %v", err)
		}
		log.Println("Migraciones realizadas exitosamente")
	}

	return db, nil
}

func InitLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
}
