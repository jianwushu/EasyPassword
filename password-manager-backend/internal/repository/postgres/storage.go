package postgres

import (
	"context"
	"easy-password-backend/internal/core"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Storage 为 PostgreSQL 实现了 repository.Storage 接口。
type Storage struct {
	db *gorm.DB
}

// NewPostgresStorage 创建一个新的 PostgreSQL 存储实例。
func NewPostgresStorage(db *gorm.DB) *Storage {
	return &Storage{db: db}
}

// User 返回一个在 PostgreSQL 数据库上操作的 UserRepository。
func (s *Storage) User() core.UserRepository {
	return &userRepository{db: s.db}
}

// Vault 返回一个在 PostgreSQL 数据库上操作的 VaultRepository。
func (s *Storage) Vault() core.VaultRepository {
	return &vaultRepository{db: s.db}
}

// --- 用户存储库实现 ---

type userRepository struct {
	db *gorm.DB
}

func (r *userRepository) Create(ctx context.Context, user *core.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*core.User, error) {
	var user core.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, core.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// --- 保险库存储库实现 ---

type vaultRepository struct {
	db *gorm.DB
}

func (r *vaultRepository) Create(ctx context.Context, item *core.VaultItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *vaultRepository) FindByID(ctx context.Context, id uuid.UUID) (*core.VaultItem, error) {
	var item core.VaultItem
	err := r.db.WithContext(ctx).First(&item, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, core.ErrVaultItemNotFound
		}
		return nil, err
	}
	return &item, nil
}

func (r *vaultRepository) FindByUser(ctx context.Context, userID uuid.UUID) ([]core.VaultItem, error) {
	var items []core.VaultItem
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&items).Error
	return items, err
}

func (r *vaultRepository) Update(ctx context.Context, item *core.VaultItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *vaultRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&core.VaultItem{}, id).Error
}