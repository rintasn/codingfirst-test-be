// main.go
package main

import (
	"log"
	"net/http"
	"time"

	"main/config"
	"main/handlers"
	"main/middleware"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Inisialisasi database
	config.InitDatabase()

	// Inisialisasi router
	router := mux.NewRouter()

	// Rute publik (tidak perlu autentikasi)
	router.HandleFunc("/api/auth/register", handlers.RegisterHandler).Methods("POST")
	router.HandleFunc("/api/auth/login", handlers.LoginHandler).Methods("POST")

	// Rute untuk manajemen preferensi (memerlukan autentikasi)
	protectedRouter := router.PathPrefix("/api").Subrouter()
	protectedRouter.Use(middleware.AuthMiddleware)

	protectedRouter.HandleFunc("/preferences", handlers.GetPreferencesHandler).Methods("GET")
	protectedRouter.HandleFunc("/preferences", handlers.UpdatePreferencesHandler).Methods("POST")
	protectedRouter.HandleFunc("/user", handlers.GetUserHandler).Methods("GET")

	// Rute untuk Claude Desktop (memerlukan autentikasi)
	protectedRouter.HandleFunc("/claude", handlers.ClaudeHandler).Methods("POST")

	// Konfigurasi CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Sesuaikan untuk produksi
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		ExposedHeaders:   []string{""},
		AllowCredentials: true,
		MaxAge:           int(12 * time.Hour / time.Second),
	})

	// Bungkus router dengan middleware CORS
	handler := c.Handler(router)

	// Mulai server
	serverAddr := ":8080"
	log.Printf("Server starting on %s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, handler))
}
