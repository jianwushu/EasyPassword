package boltdb

import (
	"easy-password-backend/internal/core"

	"go.etcd.io/bbolt"
)

var (
	userBucket             = []byte("users")
	vaultBucket            = []byte("vaults")
	usernameBucket         = []byte("usernames")
	emailBucket            = []byte("emails")
	verificationCodeBucket = []byte("verification_codes")
)

// Storage 为 BoltDB 实现了 repository.Storage 接口。
type Storage struct {
	db *bbolt.DB
}

// NewBoltDBStorage 创建一个新的 BoltDB 存储实例。
func NewBoltDBStorage(db *bbolt.DB) *Storage {
	return &Storage{db: db}
}

// User 返回一个在 BoltDB 数据库上操作的 UserRepository。
func (s *Storage) User() core.UserRepository {
	return &userRepository{db: s.db}
}

// Vault 返回一个在 BoltDB 数据库上操作的 VaultRepository。
func (s *Storage) Vault() core.VaultRepository {
	return &vaultRepository{db: s.db}
}

// VerificationCode 返回一个在 BoltDB 数据库上操作的 VerificationCodeRepository。
func (s *Storage) VerificationCode() core.VerificationCodeRepository {
	return &verificationCodeRepository{db: s.db}
}

 