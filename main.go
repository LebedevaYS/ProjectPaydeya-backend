package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {
	// Health check endpoint
	http.HandleFunc("/health", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
	}))

	// API endpoints
	http.HandleFunc("/api/", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{
			"message":   "Paydeya API is running! 🚀",
			"timestamp": time.Now().Format(time.RFC3339),
			"version":   "1.0.0",
		}
		fmt.Fprintf(w, `{"message":"%s","timestamp":"%s","version":"%s"}`,
			response["message"], response["timestamp"], response["version"])
	}))

	// Root endpoint
	http.HandleFunc("/", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"Welcome to Paydeya API","endpoints":["/health","/api/"]}`))
	}))

	// Get port from environment variable (Render sets this automatically)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("🚀 Paydeya backend server starting on port %s", port)
	log.Printf("📊 Health check: http://localhost:%s/health", port)
	log.Printf("🔗 API: http://localhost:%s/api/", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}