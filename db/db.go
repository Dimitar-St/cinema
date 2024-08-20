package db

import (
	"cinema/db/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=localhost user=postgres password=12345 dbname=postgres port=5432",
		PreferSimpleProtocol: true, // disables implicit prepared statement usage

	}), &gorm.Config{})
	if err != nil {
		println(err)
	}

	db.AutoMigrate(&models.Movie{})
	db.AutoMigrate(&models.User{})
}
