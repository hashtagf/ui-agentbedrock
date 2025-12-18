package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/ui-agentbedrock/backend/internal/config"
	"github.com/ui-agentbedrock/backend/internal/handlers"
	"github.com/ui-agentbedrock/backend/internal/repository"
	"github.com/ui-agentbedrock/backend/internal/services"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load config
	cfg := config.Load()

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoDBURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	// Ping MongoDB
	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	log.Println("Connected to MongoDB")

	db := mongoClient.Database(cfg.DatabaseName)

	// Initialize repositories
	sessionRepo := repository.NewSessionRepository(db)
	documentRepo := repository.NewDocumentRepository(db)

	// Initialize services
	sessionService := services.NewSessionService(sessionRepo)
	agentService, err := services.NewAgentService(cfg.AgentID, cfg.AgentAliasID, cfg.AgentName, cfg.AWSRegion)
	if err != nil {
		log.Fatalf("Failed to initialize agent service: %v", err)
	}
	summarizeService := services.NewSummarizeService(agentService.GetAWSConfig())
	extractService := services.NewExtractionService()

	// Initialize handlers
	sessionHandler := handlers.NewSessionHandler(sessionService)
	chatHandler := handlers.NewChatHandler(agentService, sessionService, summarizeService, documentRepo)
	uploadHandler := handlers.NewUploadHandler(documentRepo, extractService)

	// Setup Gin router
	r := gin.Default()

	// CORS configuration
	origins := strings.Split(cfg.AllowedOrigins, ",")
	r.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := r.Group("/api")
	{
		// Session routes
		api.GET("/sessions", sessionHandler.GetSessions)
		api.POST("/sessions", sessionHandler.CreateSession)
		api.GET("/sessions/:id", sessionHandler.GetSession)
		api.PUT("/sessions/:id", sessionHandler.UpdateSession)
		api.DELETE("/sessions/:id", sessionHandler.DeleteSession)
		api.DELETE("/sessions/:id/messages", sessionHandler.ClearMessages)
		api.GET("/sessions/:id/stats", sessionHandler.GetMessageStats)

		// Chat routes
		api.POST("/chat/stream", chatHandler.StreamChat)

		// Document upload routes
		api.POST("/upload", uploadHandler.UploadFile)
		api.GET("/files/:id", uploadHandler.DownloadFile)
		api.DELETE("/files/:id", uploadHandler.DeleteFile)
		api.GET("/sessions/:id/documents", uploadHandler.GetSessionDocuments)
	}

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
