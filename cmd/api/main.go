package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"secure-service/internal/auth"
	"secure-service/internal/database"
	api "secure-service/internal/http"
	"secure-service/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	auth.InitAuth()

	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.CloseDB()

	http.HandleFunc(
		"/register",
		api.ApiCommonMiddleware(
			http.MethodPost,
			api.ModelBodyMiddleware[service.RegisterRequest](
				api.RegisterHandler,
			),
		),
	)
	http.HandleFunc(
		"/login",
		api.ApiCommonMiddleware(
			http.MethodPost,
			api.ModelBodyMiddleware[service.LoginRequest](
				api.LoginHandler,
			),
		),
	)
	http.HandleFunc(
		"/profile",
		api.ApiCommonMiddleware(
			http.MethodGet,
			api.AuthMiddleware(
				api.ProfileHandler,
			),
		),
	)
	http.HandleFunc(
		"/health",
		api.ApiCommonMiddleware(
			http.MethodGet,
			api.HealthHandler,
		),
	)

	port := service.GetEnv("SERVER_PORT", "8080")
	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📝 Register: POST http://localhost:%s/register", port)
	log.Printf("🔐 Login: POST http://localhost:%s/login", port)
	log.Printf("👤 Profile: GET http://localhost:%s/profile (requires token)", port)
	log.Printf("❤️  Health: GET http://localhost:%s/health", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
