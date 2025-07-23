package core

import (
	"context"

	"github.com/google/uuid"
)

// UserRepository 定义了用户数据操作的接口。
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByUsername(ctx context.Context, username string) (*User, error)
}

// VaultRepository 定义了保险库数据操作的接口。
type VaultRepository interface {
	Create(ctx context.Context, item *VaultItem) error
	FindByID(ctx context.Context, id uuid.UUID) (*VaultItem, error)
	FindByUser(ctx context.Context, userID uuid.UUID) ([]VaultItem, error)
	Update(ctx context.Context, item *VaultItem) error
	Delete(ctx context.Context, id uuid.UUID) error
}