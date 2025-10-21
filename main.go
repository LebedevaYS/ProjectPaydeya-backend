package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

// CORS middleware
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем все origins для разработки
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		// Обработка preflight OPTIONS запросов
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
		log.Printf("Health check requested from %s", r.RemoteAddr)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "ok",
			"timestamp": time.Now().Format(time.RFC3339),
			"service":   "paydeya-backend",
		})
	}))

	// API endpoints
	http.HandleFunc("/api/", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("API request: %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":   "Paydeya API is running! 🚀",
			"timestamp": time.Now().Format(time.RFC3339),
			"endpoint":  r.URL.Path,
		})
	}))

	// Root endpoint
	http.HandleFunc("/", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Root request: %s", r.URL.Path)

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
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📊 Endpoints:")
	log.Printf("   GET /")
	log.Printf("   GET /health")
	log.Printf("   GET /api/*")
	log.Printf("🌐 CORS enabled for all origins")

	log.Fatal(http.ListenAndServe(":"+port, nil))
}