package main

import (
	"log"
	"net/http"

	"heal-her-backend-2/cmd/internal/auth"
	"heal-her-backend-2/cmd/internal/config"
	"heal-her-backend-2/cmd/internal/database"
	"heal-her-backend-2/cmd/internal/kids/dashboard"
)

func main() {
	log.Println("[BOOT] Initializing HEAL Her Go Microservice...")

	// 1. Load Configuration
	cfg := config.LoadConfig()

	// 2. Initialize the Live Database Connection Pool
	db := database.Connect(cfg)

	// 3. Initialize Live Dependencies
	// We inject the live 'db' connection into our new Postgres repository
	liveRepo := dashboard.NewPostgresDashboardRepo(db)
	clock := dashboard.DefaultTimeProvider{}

	// Inject the live repo into the service
	greetingService := dashboard.NewDashboardService(liveRepo, clock)

	// Inject the service into the handler
	kidsDashboardHandler := dashboard.NewHandler(greetingService)

	// 4. Multiplexer & Middleware setup
	mux := http.NewServeMux()
	requireAuth := auth.Middleware(cfg)

	// 5. Mount Secure Route
	mux.Handle("/kids/greet", requireAuth(http.HandlerFunc(kidsDashboardHandler.Greet)))

	// 6. Start Server
	log.Printf("[SERVER] Live Architecture locked. Listening on port %s...\n", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatalf("[FATAL] Server encountered a critical failure: %v\n", err)
	}
}
