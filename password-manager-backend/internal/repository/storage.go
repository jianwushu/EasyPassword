package repository

import (
	"easy-password-backend/config"
	"easy-password-backend/internal/core"
	"easy-password-backend/internal/repository/boltdb"
	"easy-password-backend/internal/repository/postgres"
	"fmt"

	"go.etcd.io/bbolt"
	"gorm.io/gorm"
)

// Storage 定义了数据库操作的通用接口，
// 作为特定存储库的工厂。
type Storage interface {
	User() core.UserRepository
	Vault() core.VaultRepository
}

// NewStorage 根据提供的配置创建一个新的存储后端。
// 它充当工厂并返回适当的实现（Postgres 或 BoltDB）。
func NewStorage(cfg *config.Config, db *gorm.DB, boltDB *bbolt.DB) (Storage, error) {
	switch cfg.DBType {
	case "postgres":
		return postgres.NewPostgresStorage(db), nil
	case "boltdb":
		return boltdb.NewBoltDBStorage(boltDB), nil
	default:
		return nil, fmt.Errorf("unsupported DB_TYPE: %s", cfg.DBType)
	}
}