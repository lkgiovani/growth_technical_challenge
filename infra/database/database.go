package database

import (
	"embed"
	"log"
	"os"
	"sort"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

var DB *gorm.DB

func Connect() error {
	dsn := getDSN()

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
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

	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return err
	}

	var migrationFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			migrationFiles = append(migrationFiles, entry.Name())
		}
	}

	sort.Strings(migrationFiles)

	for _, filename := range migrationFiles {
		log.Printf("Running migration: %s", filename)

		content, err := migrationsFS.ReadFile("migrations/" + filename)
		if err != nil {
			log.Printf("Error reading migration file %s: %v", filename, err)
			continue
		}

		statements := strings.Split(string(content), ";")
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}

			if err := DB.Exec(stmt).Error; err != nil {
				if strings.Contains(err.Error(), "already exists") ||
					strings.Contains(err.Error(), "duplicate key") ||
					strings.Contains(err.Error(), "relation") && strings.Contains(err.Error(), "already exists") {
					log.Printf("Migration %s already applied (skipping): %v", filename, err)
					continue
				}
				log.Printf("Warning: Error executing migration %s: %v", filename, err)
			}
		}
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
