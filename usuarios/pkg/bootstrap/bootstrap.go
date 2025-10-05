package bootstrap

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kramirez/usuarios-service/internal/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró archivo .env, usando variables de entorno del sistema")
	}
}

func DBConnection() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if os.Getenv("DATABASE_DEBUG") == "true" {
		db = db.Debug()
	}

	// Auto-migrate si DATABASE_MIGRATE está en "up"
	if os.Getenv("DATABASE_MIGRATE") == "up" {
		if err := db.AutoMigrate(&domain.Usuario{}); err != nil {
			return nil, fmt.Errorf("error al realizar migraciones: %v", err)
		}
		log.Println("Migrations completed successfully")
	}

	return db, nil
}

func InitLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
}
