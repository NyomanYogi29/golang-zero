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
)

func RunServer() {
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

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Tangani jika ada path aneh yang tidak terdaftar
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true, "message": "Server API berjalan dengan baik!"}`))
	})

	mux.HandleFunc("/api/v1/register", userHandler.Register)
	mux.HandleFunc("/api/v1/login", userHandler.Login)

	mux.HandleFunc("/api/v1/logout", authMiddleware.Authenticate(userHandler.Logout))

	fmt.Println("Server is running on port :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}

}

func main() {
	RunServer()
}
