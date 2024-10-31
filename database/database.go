package database

import (
	"ecommerce-api/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


var (
	Database *gorm.DB
)

func Connect() error {
	config := config.Config()
	dsn := config.DB_CONNECTION_STRING
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	Database = db
	if err != nil {
		panic("failed to connect database, error: " + err.Error())
	}
	return nil
}

