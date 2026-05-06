package dashboard

import (
	"fmt"
	"strings"
	"time"

	"heal-her-backend-2/cmd/pkg/security"
)

// ------------------------------------------------------------------------
// 1. ADVANCED DEPENDENCY INTERFACES
// ------------------------------------------------------------------------

// Repository defines the contract for dashboard data access.
type Repository interface {
	GetFirstName(userID string) (string, error)
	GetUserStreak(userID string) (int, error)
}

// TimeProvider abstracts the system clock for brutal, deterministic testing.
type TimeProvider interface {
	Now() time.Time
}

// DefaultTimeProvider implements TimeProvider for production use.
type DefaultTimeProvider struct{}

func (d DefaultTimeProvider) Now() time.Time {
	return time.Now()
}

// ------------------------------------------------------------------------
// 2. THE SERVICE IMPLEMENTATION
// ------------------------------------------------------------------------

// dashboardService implements the DashboardService interface.
type dashboardService struct {
	repo  Repository
	clock TimeProvider
}

// NewDashboardService is the strict constructor.
func NewDashboardService(repo Repository, clock TimeProvider) *dashboardService {
	if repo == nil {
		panic("[CRITICAL] DashboardService requires a non-nil Repository")
	}
	if clock == nil {
		panic("[CRITICAL] DashboardService requires a non-nil TimeProvider")
	}
	return &dashboardService{
		repo:  repo,
		clock: clock,
	}
}

// ------------------------------------------------------------------------
// 3. THE BUSINESS LOGIC (PURE & IMMUTABLE)
// ------------------------------------------------------------------------

// GetDashboardGreeting formats a time-aware, gamified greeting using userID.
// It returns TWO distinct strings: the full sentence and the raw sanitized name.
func (s *dashboardService) GetDashboardGreeting(userID string) (string, string, error) {
	// Step 1: Fetch user's first name from database
	firstName, err := s.repo.GetFirstName(userID)
	if err != nil {
		// Fallback to 'Friend' if DB error
		firstName = "Friend"
	}

	// Sanitize the name for security
	safeName := security.SanitizeName(firstName)
	if safeName == "" {
		safeName = "Friend"
	}

	// Step 2: Temporal Logic Calculation
	currentTime := s.clock.Now()
	hour := currentTime.Hour()
	var timeGreeting string

	switch {
	case hour >= 5 && hour < 12:
		timeGreeting = "Good morning"
	case hour >= 12 && hour < 17:
		timeGreeting = "Good afternoon"
	case hour >= 17 && hour < 20:
		timeGreeting = "Good evening"
	default:
		// Nighttime logic - gentle prompt for kids
		timeGreeting = "It's getting late"
	}

	// Step 3: Gamification Fetch (Currently bypassed to return 0)
	streak, err := s.repo.GetUserStreak(userID)
	if err != nil {
		fullGreeting := fmt.Sprintf("%s, %s! Welcome to your safe space.", timeGreeting, safeName)
		return fullGreeting, safeName, nil
	}

	// Step 4: High-Performance String Construction
	var builder strings.Builder
	builder.WriteString(timeGreeting)
	builder.WriteString(", ")
	builder.WriteString(safeName)
	builder.WriteString("!")

	if streak > 1 {
		builder.WriteString(fmt.Sprintf(" You're on a %d-day streak! Keep it up!", streak))
	} else {
		builder.WriteString(" Welcome to your safe space.")
	}

	return builder.String(), safeName, nil
}