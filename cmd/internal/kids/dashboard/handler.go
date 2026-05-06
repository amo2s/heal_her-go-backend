package dashboard

import (
	"encoding/json"
	"log"
	"net/http"

	"heal-her-backend-2/cmd/internal/auth" 
)

// ------------------------------------------------------------------------
// 1. THE SERVICE INTERFACE
// ------------------------------------------------------------------------
// The interface now requires the service to return TWO strings:
// 1. The full temporal greeting (for the dashboard hero)
// 2. The sanitized first name (for the TopBar)
type DashboardService interface {
	GetDashboardGreeting(userID string) (string, string, error)
}

// Handler holds the dependencies for the dashboard HTTP routes.
type Handler struct {
	service DashboardService
}

// NewHandler is the constructor function.
func NewHandler(svc DashboardService) *Handler {
	return &Handler{
		service: svc,
	}
}

// ------------------------------------------------------------------------
// 2. RESPONSE SCHEMAS
// ------------------------------------------------------------------------
// We updated the Data struct to hold both specific string targets.
type GreetingResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		FullGreeting string `json:"full_greeting,omitempty"`
		FirstName    string `json:"first_name,omitempty"`
	} `json:"data,omitempty"`
}

// ------------------------------------------------------------------------
// 3. THE HTTP METHOD LOGIC
// ------------------------------------------------------------------------

// Greet handles the incoming HTTP request for the kids' dashboard greeting.
func (h *Handler) Greet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Step 1: Extract Identity from Context
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.UserClaims)
	if !ok || claims == nil {
		log.Println("[CRITICAL] Dashboard handler accessed without Auth Context")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(GreetingResponse{
			Status:  "error",
			Message: "Internal server architecture error",
		})
		return
	}

	// Step 2: Hard-Boundary Segment Validation (The Python Lock)
	if claims.Dashboard != "kids" {
		log.Printf("[SECURITY] User ID %s assigned to '%s' blocked from Kids segment\n", claims.ID, claims.Dashboard)
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(GreetingResponse{
			Status:  "error",
			Message: "Access denied. Segment boundary violation.",
		})
		return
	}

	// Step 3: Execute Business Logic via the ID
	// We now capture both strings returned by the updated service
	fullGreeting, firstName, err := h.service.GetDashboardGreeting(claims.ID)
	if err != nil {
		log.Printf("[ERROR] Service failed to generate greeting for ID %s: %v\n", claims.ID, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(GreetingResponse{
			Status:  "error",
			Message: "Unable to load dashboard data at this time.",
		})
		return
	}

	// Step 4: Dispatch Response
	// Pack both boxes into the JSON delivery
	response := GreetingResponse{
		Status:  "success",
		Message: "Dashboard loaded successfully",
	}
	response.Data.FullGreeting = fullGreeting
	response.Data.FirstName = firstName

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}