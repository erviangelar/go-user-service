package db

import (
	"log"

	"github.com/erviangelar/go-users-api/pkg/common/config"
	"github.com/erviangelar/go-users-api/pkg/common/models"
	"github.com/jmoiron/sqlx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// func Init(url string) *gorm.DB {
// 	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})

// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	db.AutoMigrate(&models.User{})

//		return db
//	}

func Init(config *config.Configurations) {
	db, err := gorm.Open(postgres.Open(config.DBConn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&models.User{})
}

// Create Connection to database using sqlx
func NewConnection(config *config.Configurations) (*sqlx.DB, error) {
	var conn string

	if config.DBConn != "" {
		conn = config.DBConn
	}

	db, err := sqlx.Connect("postgres", conn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
