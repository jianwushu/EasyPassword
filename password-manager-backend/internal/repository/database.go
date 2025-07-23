package repository

import (
	"easy-password-backend/config"
	"fmt"
	"time"

	"go.etcd.io/bbolt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect 初始化 PostgreSQL 数据库连接。
func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.DatabaseURL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}

// InitBoltDB 初始化 BoltDB 数据库并创建必要的存储桶。
func InitBoltDB(path string) (*bbolt.DB, error) {
	db, err := bbolt.Open(path, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		buckets := [][]byte{
			[]byte("users"),
			[]byte("vaults"),
			[]byte("usernames"),
		}
		for _, bucket := range buckets {
			_, err := tx.CreateBucketIfNotExists(bucket)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return db, nil
}