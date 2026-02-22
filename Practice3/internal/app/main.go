package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	"golang/internal/handlers"
	"golang/internal/middleware"
	"golang/internal/repository"
	_postgres "golang/internal/repository/_postgres"
	userUsecase "golang/internal/usecase/users"
	"golang/pkg/modules"
)

func Run() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConfig := initPostgreConfig()
	serverConfig := initServerConfig()
	redisConfig := initRedisConfig()

	pgDialect := _postgres.NewPGXDialect(ctx, dbConfig)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Printf("Redis connection failed: %v (continuing without cache)", err)
		redisClient = nil
	} else {
		log.Println("Connected to Redis")
	}

	repos := repository.NewRepositories(pgDialect)
	userUC := userUsecase.NewUserUsecase(repos.UserRepository, redisClient)

	userHandler := handlers.NewUserHandler(userUC)
	authHandler := handlers.NewAuthHandler(userUC, serverConfig.JWTSecret)

	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logging)
	r.Use(middleware.APIKeyAuth(serverConfig.APIKey))

	// Public routes (API key only)
	r.Get("/health", handlers.HealthCheck)
	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)

	// Swagger (API key only)
	r.Get("/swagger/*", handlers.SwaggerHandler())

	// Protected routes (API key + JWT)
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth(serverConfig.JWTSecret))

		r.Get("/users", userHandler.GetUsers)
		r.Get("/users/{id}", userHandler.GetUserByID)
		r.Post("/users", userHandler.CreateUser)
		r.Put("/users/{id}", userHandler.UpdateUser)
		r.Post("/users/with-audit", userHandler.CreateUserWithAudit)

		// Admin only
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole("admin"))
			r.Delete("/users/{id}", userHandler.DeleteUser)
		})
	})

	// Background worker: count users every 60 seconds
	ticker := time.NewTicker(60 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				count, err := userUC.CountUsers()
				if err != nil {
					log.Printf("[Background Worker] Error counting users: %v", err)
				} else {
					log.Printf("[Background Worker] Total active users: %d", count)
				}
			case <-ctx.Done():
				ticker.Stop()
				log.Println("[Background Worker] Stopped")
				return
			}
		}
	}()

	srv := &http.Server{
		Addr:    ":" + serverConfig.Port,
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigChan
		log.Printf("Received signal: %v. Shutting down...", sig)

		cancel()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}

		if redisClient != nil {
			if err := redisClient.Close(); err != nil {
				log.Printf("Redis close error: %v", err)
			}
		}

		if err := pgDialect.DB.Close(); err != nil {
			log.Printf("Database close error: %v", err)
		}

		log.Println("Server stopped gracefully")
	}()

	log.Printf("Server starting on :%s", serverConfig.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}

func initPostgreConfig() *modules.PostgreConfig {
	timeout, _ := time.ParseDuration(getEnv("DB_EXEC_TIMEOUT", "5s"))
	return &modules.PostgreConfig{
		Host:        getEnv("DB_HOST", "localhost"),
		Port:        getEnv("DB_PORT", "5432"),
		Username:    getEnv("DB_USERNAME", "postgres"),
		Password:    getEnv("DB_PASSWORD", "postgres"),
		DBName:      getEnv("DB_NAME", "mydb"),
		SSLMode:     getEnv("DB_SSLMODE", "disable"),
		ExecTimeout: timeout,
	}
}

func initServerConfig() *modules.ServerConfig {
	return &modules.ServerConfig{
		Port:      getEnv("SERVER_PORT", "8080"),
		APIKey:    getEnv("API_KEY", "secret12345"),
		JWTSecret: getEnv("JWT_SECRET", "my-super-secret-jwt-key"),
	}
}

func initRedisConfig() *modules.RedisConfig {
	db, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	return &modules.RedisConfig{
		Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       db,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
