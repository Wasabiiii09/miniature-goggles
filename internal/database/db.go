package database

import (
	"log"
	"os"

	"work/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatalln("Link zu DB nicht gefunden in file .env")
	}

	var err error 
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Fehler beim Verbindung zur DatenBank: %v", err)
	}

	err = DB.AutoMigrate(&models.User{}, &models.Comment{})
	if err != nil {
		log.Fatalf("Fehler beim AutoMigrate: %v", err)
	}

	log.Println("Datenbank erfolgreich verbunden und migriert!")

}