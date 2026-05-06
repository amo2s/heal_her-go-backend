package database

import (
	"log"
	"os"
	"time"

	"heal-her-backend-2/cmd/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global access point for the database connection.
var DB *gorm.DB

// Connect initializes the PostgreSQL connection pool.
func Connect(cfg *config.Config) *gorm.DB {
	// 1. Validation Check
	if cfg.DatabaseURL == "" {
		log.Fatal("[FATAL] Database connection failed: DATABASE_URL is empty in configuration.")
	}

	// 2. Custom Logger for High-Latency Environments
	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second, // 1000ms threshold to silence normal geographical latency
			LogLevel:                  logger.Warn, // Only log warnings and errors
			IgnoreRecordNotFoundError: true,        // Ignore standard 404s in logs
			Colorful:                  true,        // Maintain console formatting
		},
	)

	// 3. Open Connection with strict PgBouncer (Pooler) compatibility
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  cfg.DatabaseURL,
		PreferSimpleProtocol: true, // CRITICAL FIX 1: Forces the pgx driver to bypass prepared statements entirely
	}), &gorm.Config{
		PrepareStmt: false,    // CRITICAL FIX 2: Prevents GORM from caching statements at a higher level
		Logger:      dbLogger, // Apply the custom latency threshold
	})

	if err != nil {
		log.Fatalf("[FATAL] Failed to connect to database: %v", err)
	}

	// 4. Brutal Connection Pool Tuning
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("[FATAL] Failed to access underlying SQL driver: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
	log.Println("[DATABASE] Connection established. Driver locked to simple protocol for pooler compatibility.")
	return db
}

// User is the strict mapping of the SQLAlchemy 'users' table.
// WARNING: Do not run db.AutoMigrate(&User{}) in Go. Python owns the schema.
type User struct {
	ID           string    `gorm:"primaryKey;column:id"`
	FullName     string    `gorm:"column:full_name;not null"`
	Email        string    `gorm:"uniqueIndex;column:email;not null"`
	Age          int       `gorm:"column:age;not null"`
	PasswordHash string    `gorm:"column:password_hash;not null"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

// TableName explicitly binds this struct to the existing "users" table.
func (User) TableName() string {
	return "users"
}