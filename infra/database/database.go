package database

import (
	"log"
	"os"

	"github.com/lkgiovani/growth_technical_challenge/internal/domain/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() error {
	dsn := getDSN()

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		return err
	}

	log.Println("Database connected successfully")

	if err := AutoMigrate(); err != nil {
		return err
	}

	return nil
}

func AutoMigrate() error {
	log.Println("Running auto migrations...")

	if err := DB.AutoMigrate(
		&entities.Departamento{},
		&entities.Colaborador{},
	); err != nil {
		return err
	}

	log.Println("Auto migrations completed successfully")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}

func getDSN() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("POSTGRES_USER", "postgres")
	password := getEnv("POSTGRES_PASSWORD", "postgres")
	dbname := getEnv("POSTGRES_DB", "growth_db")
	sslmode := getEnv("DB_SSLMODE", "disable")

	return "host=" + host + " port=" + port + " user=" + user +
		" password=" + password + " dbname=" + dbname + " sslmode=" + sslmode
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
