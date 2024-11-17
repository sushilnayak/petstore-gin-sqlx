package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"petstore/config"
	"petstore/internal/handler"
	"petstore/internal/middleware"
	"petstore/internal/repository"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database connection
	dbURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	petRepo := repository.NewPetRepository(db)

	// Initialize handlers
	petHandler := handler.NewPetHandler(petRepo)

	// Setup Gin router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())

	// API routes
	v1 := router.Group("/v1")
	{
		// Public routes
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.APIKeyAuth(cfg.Auth.APIKey))
		{
			// Pet routes
			protected.POST("/pet", petHandler.CreatePet)
			protected.GET("/pet/:petId", petHandler.GetPet)
			protected.PUT("/pet", petHandler.UpdatePet)
			protected.DELETE("/pet/:petId", petHandler.DeletePet)
		}
	}

	// Start server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := router.Run(cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
