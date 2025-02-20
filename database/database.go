package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go_api_example/models"
)

var DB *gorm.DB

func DropUnusedColumns(db *gorm.DB, dst ...interface{}) {
	for _, m := range dst {
		stmt := &gorm.Statement{DB: db}
		stmt.Parse(m)
		fields := stmt.Schema.Fields
		columns, _ := db.Migrator().ColumnTypes(m)

		for i := range columns {
			found := false
			for j := range fields {
				if columns[i].Name() == fields[j].DBName {
					found = true
					break
				}
			}
			if !found {
				db.Migrator().DropColumn(m, columns[i].Name())
			}
		}
	}
}

// Auto-migrate tables
func MigrateDB(db *gorm.DB, dst ...interface{}) {
	err := db.AutoMigrate(dst...)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	DropUnusedColumns(db, dst...)
	log.Println("Database migration completed!")
}

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = db
	MigrateDB(DB, &models.Product{}, &models.User{})
	log.Println("Database connection successful!")
}
