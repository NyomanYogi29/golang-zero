package main

import (
	"context"
	"fmt"
	"latihan/internal/config"
	"latihan/internal/database"
	"latihan/internal/handler"
	"latihan/internal/middleware"
	"latihan/internal/repository"
	"latihan/internal/service"
	"log"
	"net/http"
	"os"

	_ "latihan/docs"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/joho/godotenv"
)

// @title Go-Playground API Documentation
// @version 1.0
// @description This is the main entry point of the server
// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token.

func RunServer() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file is not found. Fallback into default OS environment")
	}

	if err := config.ConnectPostgres(); err != nil {
		log.Fatal(err)
	}

	if err := config.ConnectRedis(); err != nil {
		log.Fatal(err)
	}

	dbPool := config.GetDBPool()
	rdb := config.GetRedisClient()

	if err := database.RunMigrations(context.Background(), dbPool); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	userRepo := repository.NewPostgresUserRepository(dbPool)
	userService := service.NewUserService(userRepo, rdb)
	userHandler := handler.NewUserService(userService)
	authMiddleware := middleware.NewAuthMiddleware(userService)

	mux := http.NewServeMux()

	mux.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true, "message": "Gateway is running like a butter!"}`))
	})

	mux.HandleFunc("/api/v1/register", userHandler.Register)
	mux.HandleFunc("/api/v1/login", userHandler.Login)

	mux.HandleFunc("/api/v1/logout", authMiddleware.Authenticate(userHandler.Logout))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server is running on port :%s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}

}

func main() {
	RunServer()
}
