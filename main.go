package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// Health check endpoint - ДОЛЖЕН БЫТЬ ПЕРЕД корневым маршрутом
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Health check requested from %s", r.RemoteAddr)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "ok",
			"timestamp": time.Now().Format(time.RFC3339),
			"service":   "paydeya-backend",
		})
	})

	// API endpoints
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("API request: %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":   "Paydeya API is running! 🚀",
			"timestamp": time.Now().Format(time.RFC3339),
			"endpoint":  r.URL.Path,
		})
	})

	// Root endpoint - ДОЛЖЕН БЫТЬ ПОСЛЕДНИМ
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Root request: %s", r.URL.Path)

		// Если запрос не к корню, покажем 404
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Welcome to Paydeya API",
			"endpoints": []string{
				"GET /health",
				"GET /api/",
			},
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📊 Endpoints:")
	log.Printf("   GET /")
	log.Printf("   GET /health")
	log.Printf("   GET /api/*")

	log.Fatal(http.ListenAndServe(":"+port, nil))
}