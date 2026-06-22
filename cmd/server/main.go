package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"recruitment-platform/internal/application"
	"recruitment-platform/internal/auth"
	"recruitment-platform/internal/chat"
	"recruitment-platform/internal/config"
	"recruitment-platform/internal/interview"
	"recruitment-platform/internal/job"
	"recruitment-platform/internal/recommend"
	"recruitment-platform/internal/middleware"
	"recruitment-platform/internal/resume"
	"recruitment-platform/internal/user"
	pkgauth "recruitment-platform/internal/pkg/auth"
	"recruitment-platform/internal/pkg/broker"
	"recruitment-platform/internal/pkg/cos"
	"recruitment-platform/internal/pkg/maps"
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

	authRepo := auth.NewRepository(dbPool)
	authSvc := auth.NewService(authRepo, jwtManager, smsSender, redisClient, msgBroker)
	authHandler := auth.NewHandler(authSvc)

	// ── 初始化 COS Uploader ──
	var cosUploader cos.Uploader
	if cfg.COS.SecretID != "" {
		cosUploader, err = cos.NewTencentCOS(cfg.COS)
		if err != nil {
			log.Fatalf("Failed to init COS client: %v", err)
		}
		log.Println("COS client initialized")
	} else {
		log.Println("COS client: not configured (attachment upload will fail)")
	}

	// ── 初始化 User 模块 ──
	userRepo := user.NewRepository(dbPool)
	userSvc := user.NewService(userRepo)
	userHandler := user.NewHandler(userSvc)

	// ── 初始化 Resume 模块 ──
	resumeRepo := resume.NewRepository(dbPool)
	resumeSvc := resume.NewService(resumeRepo, cosUploader)
	resumeHandler := resume.NewHandler(resumeSvc)

	// ── 初始化 Job 模块 ──
	jobRepo := job.NewRepository(dbPool)
	jobSvc := job.NewService(jobRepo)
	jobHandler := job.NewHandler(jobSvc)

	// ── 初始化 Application 模块 ──
	appRepo := application.NewRepository(dbPool)
	appSvc := application.NewService(appRepo)
	appHandler := application.NewHandler(appSvc)

	// ── 初始化 Chat/IM 模块 ──
	chatHandler := chat.NewHandler(cfg.IM.AppID, cfg.IM.Secret)

	// ── 初始化 Maps 客户端 ──
	mapsClient := maps.NewClient(cfg.Maps.APIKey)

	// ── 初始化 Interview 模块 ──
	interviewRepo := interview.NewRepository(dbPool)
	interviewSvc := interview.NewService(interviewRepo, dbPool)
	interviewHandler := interview.NewHandler(interviewSvc, mapsClient)

	// ── 初始化 Recommend 模块 ──
	recommendSvc := recommend.NewService(dbPool)
	recommendHandler := recommend.NewHandler(recommendSvc)

	router := gin.New()

	router.Use(middleware.Logger())
	router.Use(middleware.CORS())
	router.Use(gin.Recovery())

	limiter := middleware.NewSimpleRateLimiter(100)
	router.Use(limiter.Middleware())

	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 注册业务路由
	api := router.Group("/api")
	authHandler.RegisterRoutes(api)

	// 需要认证的路由
	protected := api.Group("")
	protected.Use(middleware.AuthRequired(jwtManager))
	{
		userHandler.RegisterRoutes(protected)
		resumeHandler.RegisterRoutes(protected)
		jobHandler.RegisterRoutes(api, protected)
		appHandler.RegisterRoutes(api, protected)
		chatHandler.RegisterRoutes(protected)
		interviewHandler.RegisterRoutes(api, protected)
			recommendHandler.RegisterRoutes(protected)
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
