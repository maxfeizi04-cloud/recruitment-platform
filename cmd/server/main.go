package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"recruitment-platform/internal/auth"
	"recruitment-platform/internal/config"
	"recruitment-platform/internal/middleware"
	pkgauth "recruitment-platform/internal/pkg/auth"
	"recruitment-platform/internal/pkg/broker"
	"recruitment-platform/internal/pkg/cos"
	redisclient "recruitment-platform/internal/pkg/redis"
	"recruitment-platform/internal/pkg/sms"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	configPath := "config/config.yaml"
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		configPath = envPath
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	gin.SetMode(cfg.Server.Mode)

	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, cfg.Database.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()
	log.Println("Database connected")

	redisClient, err := redisclient.NewClient(cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()
	log.Println("Redis connected")

	jwtManager := pkgauth.NewJWTManager(cfg.JWT)

	var smsSender sms.Sender
	if cfg.SMS.SecretID != "" {
		smsSender, err = sms.NewTencentSMS(cfg.SMS)
		if err != nil {
			log.Fatalf("Failed to init SMS client: %v", err)
		}
		log.Println("SMS client initialized (Tencent Cloud)")
	} else {
		smsSender = &sms.MockSender{}
		log.Println("SMS client: using mock (no SMS config provided)")
	}

	msgBroker := broker.NewInMemoryBroker()
	defer msgBroker.Close()

	_ = cos.NewTencentCOS

	authRepo := auth.NewRepository(dbPool)
	authSvc := auth.NewService(authRepo, jwtManager, smsSender, redisClient, msgBroker)
	authHandler := auth.NewHandler(authSvc)

	router := gin.New()

	router.Use(middleware.Logger())
	router.Use(middleware.CORS())
	router.Use(gin.Recovery())

	limiter := middleware.NewSimpleRateLimiter(100)
	router.Use(limiter.Middleware())

	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := router.Group("/api")
	authHandler.RegisterRoutes(api)

	protected := api.Group("")
	protected.Use(middleware.AuthRequired(jwtManager))
	{
		protected.GET("/users/me", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			role, _ := c.Get("role")
			c.JSON(http.StatusOK, gin.H{
				"user_id": userID,
				"role":    role,
			})
		})
	}

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	go func() {
		log.Printf("Server starting on :%s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
