package dashboard

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"heal-her-backend-2/cmd/internal/database" 

	"gorm.io/gorm"
)

// PostgresDashboardRepo is the live implementation for fetching dashboard data.
type PostgresDashboardRepo struct {
	db *gorm.DB
}

// NewPostgresDashboardRepo enforces strict dependency injection.
// It returns the concrete struct which implicitly satisfies the Repository interface in service.go.
func NewPostgresDashboardRepo(db *gorm.DB) *PostgresDashboardRepo {
	if db == nil {
		panic("[CRITICAL] PostgresDashboardRepo initialized with a nil database pointer.")
	}
	return &PostgresDashboardRepo{db: db}
}

// GetFirstName fetches the user's full name from the 'users' table and isolates the first name.
func (r *PostgresDashboardRepo) GetFirstName(userID string) (string, error) {
	// 1. Brutal Timeout Boundary
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var user database.User

	// 2. Query the 'users' table
	err := r.db.WithContext(ctx).
		Select("full_name").
		Where("id = ?", userID).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Fallback so the UI doesn't crash
			return "Friend", nil
		}
		log.Printf("[DB ERROR] Failed to fetch user for greeting: %v\n", err)
		return "Friend", err
	}

	// 3. The Name Stripper
	// Takes "Nwaka Amos Chika" and returns "Nwaka" (index 0)
	nameParts := strings.Fields(user.FullName)
	if len(nameParts) > 0 {
		return nameParts[0], nil
	}

	return user.FullName, nil
}

// GetUserStreak returns the streak count (currently bypassed until table exists).
func (r *PostgresDashboardRepo) GetUserStreak(userID string) (int, error) {
	return 0, nil
}