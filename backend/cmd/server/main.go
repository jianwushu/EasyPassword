package main

import (
	v1 "easy-password-backend/api/v1"
	"easy-password-backend/config"
	"easy-password-backend/internal/auth"
	"easy-password-backend/internal/core"
	"easy-password-backend/internal/email"
	"easy-password-backend/internal/repository"
	"easy-password-backend/internal/service"
	"easy-password-backend/pkg/logger"

	"log/slog"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"go.etcd.io/bbolt"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化日志
	level := parseLogLevel(cfg.LogLevel)
	logger.Init(level, cfg.LogFormat, os.Stdout)

	slog.Info("Configuration loaded successfully")

	// 初始化数据库连接
	var gormDB *gorm.DB
	var boltDB *bbolt.DB
	var err error

	switch cfg.DBType {
	case "postgres":
		slog.Info("Using PostgreSQL database.")
		gormDB, err = repository.Connect(cfg)
		if err != nil {
			slog.Error("could not connect to postgres", "error", err)
			os.Exit(1)
		}
		// 自动迁移模式
		err = gormDB.AutoMigrate(&core.User{}, &core.VaultItem{}, &core.VerificationCode{})
		if err != nil {
			slog.Error("Failed to migrate database", "error", err)
			os.Exit(1)
		}
		slog.Info("Database migrated successfully.")
	case "boltdb":
		slog.Info("Using BoltDB database.")
		boltDB, err = repository.InitBoltDB(cfg.DBPath)
		if err != nil {
			slog.Error("could not initialize boltdb", "error", err)
			os.Exit(1)
		}
		defer boltDB.Close()
	default:
		slog.Error("Unsupported DB_TYPE", "db_type", cfg.DBType)
		os.Exit(1)
	}

	// 创建存储后端
	storage, err := repository.NewStorage(cfg, gormDB, boltDB)
	if err != nil {
		slog.Error("could not create storage", "error", err)
		os.Exit(1)
	}

	// 初始化服务
	emailService := email.NewSMTPEmailService(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPFrom)
	authService := auth.NewAuthService(storage.User(), storage.VerificationCode(), emailService, cfg)
	slog.Info("AuthService initialized.")
	vaultService := service.NewVaultService(storage.Vault())
	slog.Info("VaultService initialized.")

	// 初始化 Gin 路由
	gin.SetMode(gin.ReleaseMode) // 设置为生产模式
	router := gin.Default()

	// 使用日志中间件
	router.Use(v1.LoggingMiddleware())

	// 初始化处理程序
	authHandler := v1.NewAuthHandler(authService)
	authHandler.RegisterRoutes(router)

	// 受保护的路由
	vaultAPI := router.Group("/api/v1")
	vaultAPI.Use(v1.AuthMiddleware(cfg))
	{
		vaultHandler := v1.NewVaultHandler(vaultService)
		vaultHandler.RegisterRoutes(vaultAPI)
	}

	// 启动服务器
	slog.Info("Starting server", "address", ":8081")
	if err := router.Run(":8081"); err != nil {
		slog.Error("Failed to run server", "error", err)
		os.Exit(1)
	}
}

func parseLogLevel(levelStr string) slog.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
