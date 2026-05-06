package database

import (
	"log"
	"time"

	"heal-her-backend-2/cmd/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect establishes a resilient connection pool to the shared database.
func Connect(cfg *config.Config) *gorm.DB {
	// GORM Configuration
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		// Prevent GORM from pluralizing table names globally, though we explicitly set it in models
		NamingStrategy: nil,
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), gormConfig)
	if err != nil {
		log.Fatalf("[FATAL] Failed to establish database connection: %v\n", err)
	}

	// Extract standard SQL DB interface for advanced pooling configuration
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("[FATAL] Failed to extract sql.DB interface: %v\n", err)
	}

	// Advanced Pool Tuning
	sqlDB.SetMaxIdleConns(20)                  // Keep 20 connections open and ready
	sqlDB.SetMaxOpenConns(100)                 // Max simultaneous connections to prevent DB exhaustion
	sqlDB.SetConnMaxLifetime(time.Hour)        // Recycle connections to avoid stale drops
	sqlDB.SetConnMaxIdleTime(30 * time.Minute) // Close idle connections after 30 mins

	log.Println("[OK] Postgres Connection Pool Initialized.")
	return db
}
