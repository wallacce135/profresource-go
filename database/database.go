package database

import (
	"fmt"
	"log"
	"os"

	"github.com/wallacce135/profresource/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DBConnection *gorm.DB
)

func ConnectToDatabase() {

	database_host := os.Getenv("DATABASE_HOST")
	database_port := os.Getenv("DATABASE_PORT")
	database_user := os.Getenv("DATABASE_USER")
	database_passwd := os.Getenv("DATABASE_PASSWORD")
	database_name := os.Getenv("DATABASE_NAME")

	database_migrate := os.Getenv("DATABASE_MIGRATE")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb3&parseTime=true", database_user, database_passwd, database_host, database_port, database_name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	if database_migrate == "true" {
		db.AutoMigrate(&models.User{}, &models.Articles{}, &models.Comments{})
	}

	DBConnection = db
}
