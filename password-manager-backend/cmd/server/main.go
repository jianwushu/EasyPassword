package main

import (
	"easy-password-backend/api/v1"
	"easy-password-backend/config"
	"easy-password-backend/internal/auth"
	"easy-password-backend/internal/core"
	"easy-password-backend/internal/repository"
	"easy-password-backend/internal/service"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"go.etcd.io/bbolt"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库连接
	var gormDB *gorm.DB
	var boltDB *bbolt.DB
	var err error

	switch cfg.DBType {
	case "postgres":
		log.Println("Using PostgreSQL database.")
		gormDB, err = repository.Connect(cfg)
		if err != nil {
			log.Fatalf("could not connect to postgres: %v", err)
		}
		// 自动迁移模式
		err = gormDB.AutoMigrate(&core.User{}, &core.VaultItem{})
		if err != nil {
			log.Fatalf("Failed to migrate database: %v", err)
		}
		fmt.Println("Database migrated successfully.")
	case "boltdb":
		log.Println("Using BoltDB database.")
		boltDB, err = repository.InitBoltDB(cfg.DBPath)
		if err != nil {
			log.Fatalf("could not initialize boltdb: %v", err)
		}
		defer boltDB.Close()
	default:
		log.Fatalf("Unsupported DB_TYPE: %s", cfg.DBType)
	}

	// 创建存储后端
	storage, err := repository.NewStorage(cfg, gormDB, boltDB)
	if err != nil {
		log.Fatalf("could not create storage: %v", err)
	}

	// 初始化服务
	authService := auth.NewAuthService(storage.User(), cfg)
	log.Println("AuthService initialized.")
	vaultService := service.NewVaultService(storage.Vault())
	log.Println("VaultService initialized.")

	// 初始化 Gin 路由
	router := gin.Default()

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
	log.Println("Starting server on :8081")
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}