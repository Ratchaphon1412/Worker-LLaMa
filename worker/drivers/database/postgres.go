package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Ratchaphon1412/worker-llama/configs"
	"github.com/Ratchaphon1412/worker-llama/pkg/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

// connectDb
func Connect(cfg *configs.Config) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", cfg.DB_HOST, cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_NAME, cfg.DB_PORT, cfg.DB_SSL_MODE, cfg.DB_TIMEZONE)

	db_logger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // allow colors
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: db_logger,
	})
	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}

	log.Println("connected")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("running migrations")
	// Migrate the schema
	db.AutoMigrate(
		entities.Chat{},
		entities.Research{},
		entities.Thumbnail{},
	)

	DB = Dbinstance{
		Db: db,
	}
}
