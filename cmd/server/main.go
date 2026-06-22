// @title           放心招聘平台 API
// @version         2.0
// @description     放心招聘平台后端接口文档，支持求职者、招聘者、管理员三端
// @contact.name    API Support
// @contact.url     https://github.com/maxfeizi04-cloud/recruitment-platform
// @host            localhost:8080
// @BasePath        /api
// @schemes         http
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
// @description     输入 Bearer {token}

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/maxfeizi04-cloud/recruitment-platform/internal/application"
	"github.com/maxfeizi04-cloud/recruitment-platform/internal/auth"
	"github.com/maxfeizi04-cloud/recruitment-platform/internal/chat"
	"github.com/maxfeizi04-cloud/recruitment-platform/internal/config"
	"github.com/maxfeizi04-cloud/recruitment-platform/internal/interview"
	"github.com/maxfeizi04-cloud/recruitment-platform/internal/job"
	"github.com/maxfeizi04-cloud/recruitment-platform/internal/recommend"
	"github.com/maxfeizi04-cloud/recruitment-platform/internal/middleware"
	"github.com/maxfeizi04-cloud/recruitment-platform/internal/resume"
	"github.com/maxfeizi04-cloud/recruitment-platform/internal/user"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/maxfeizi04-cloud/recruitment-platform/docs"
	pkgauth "github.com/maxfeizi04-cloud/recruitment-platform/internal/pkg/auth"
	"github.com/maxfeizi04-cloud/recruitment-platform/internal/pkg/broker"
	"github.com/maxfeizi04-cloud/recruitment-platform/internal/pkg/cache"
	"github.com/maxfeizi04-cloud/recruitment-platform/internal/pkg/cos"
	"github.com/maxfeizi04-cloud/recruitment-platform/internal/pkg/maps"
	redisclient "github.com/maxfeizi04-cloud/recruitment-platform/internal/pkg/redis"
	"github.com/maxfeizi04-cloud/recruitment-platform/internal/pkg/ws"
	"github.com/maxfeizi04-cloud/recruitment-platform/internal/pkg/search"
	"github.com/maxfeizi04-cloud/recruitment-platform/internal/pkg/sms"

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
		middleware.Logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	gin.SetMode(cfg.Server.Mode)

	// Set log level based on Gin mode
	if cfg.Server.Mode == "debug" {
		middleware.SetLogLevel("debug")
	}

	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, cfg.Database.DSN())
	if err != nil {
		middleware.Logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()
	middleware.Logger.Info("database connected")

	redisClient, err := redisclient.NewClient(cfg.Redis)
	if err != nil {
		middleware.Logger.Error("failed to connect to Redis", "error", err)
		os.Exit(1)
	}
	defer redisClient.Close()
	middleware.Logger.Info("redis connected")

	jwtManager := pkgauth.NewJWTManager(cfg.JWT)

	var smsSender sms.Sender
	if cfg.SMS.SecretID != "" {
		smsSender, err = sms.NewTencentSMS(cfg.SMS)
		if err != nil {
			middleware.Logger.Error("failed to init SMS client", "error", err)
			os.Exit(1)
		}
		middleware.Logger.Info("SMS client initialized (Tencent Cloud)")
	} else {
		smsSender = &sms.MockSender{}
		middleware.Logger.Info("SMS client: using mock (no SMS config provided)")
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
			middleware.Logger.Error("failed to init COS client", "error", err)
			os.Exit(1)
		}
		middleware.Logger.Info("COS client initialized")
	} else {
		middleware.Logger.Info("COS client: not configured (attachment upload will fail)")
	}

	// ── 初始化 User 模块 ──
	userRepo := user.NewRepository(dbPool)
	userSvc := user.NewService(userRepo)
	userHandler := user.NewHandler(userSvc)

	// ── 初始化 Resume 模块 ──
	resumeRepo := resume.NewRepository(dbPool)
	resumeSvc := resume.NewService(resumeRepo, cosUploader)
	resumeHandler := resume.NewHandler(resumeSvc)

	// ── 初始化缓存 ──

		// ── 初始化 Job 模块 ──
		cacheClient := cache.New(redisClient.Client)

		// ── 初始化 ES 客户端（可选）──
		var esClient *search.Client
		if cfg.ES.Addr != "" {
			esClient, err = search.NewClient([]string{cfg.ES.Addr}, cfg.ES.Index)
			if err != nil {
				middleware.Logger.Warn("ES not available, using PostgreSQL search", "error", err)
			} else {
				middleware.Logger.Info("Elasticsearch connected", "addr", cfg.ES.Addr, "index", cfg.ES.Index)
			}
		}
	jobRepo := job.NewRepository(dbPool)
	jobSvc := job.NewService(jobRepo, cacheClient, esClient)
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

	router.Use(middleware.Tracing())
	router.Use(middleware.RequestLogger())
	router.Use(middleware.CORS())
	router.Use(middleware.CharsetUTF8())
	router.Use(gin.Recovery())

	limiter := middleware.NewSimpleRateLimiter(100)
	router.Use(limiter.Middleware())

	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Swagger 文档
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// WebSocket
	ws.JWTSecret = []byte(cfg.JWT.Secret)
	go ws.DefaultHub.Run()
	router.GET("/ws", ws.HandleWebSocket)

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
		middleware.Logger.Info("server starting", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			middleware.Logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	middleware.Logger.Info("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		middleware.Logger.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	middleware.Logger.Info("server exited")
}
