package service

import (
	"context"
	"easy-password-backend/internal/apierror"
	"easy-password-backend/internal/core"
	"time"

	"github.com/google/uuid"
)

// VaultService 提供与保险库相关的服务。
type VaultService struct {
	vaultRepo core.VaultRepository
}

// NewVaultService 创建一个新的 VaultService。
func NewVaultService(vaultRepo core.VaultRepository) *VaultService {
	return &VaultService{vaultRepo: vaultRepo}
}

// CreateVaultItem 为用户创建一个新的保险库项目。
func (s *VaultService) CreateVaultItem(ctx context.Context, item *core.VaultItem) (*core.VaultItem, error) {
	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now
	err := s.vaultRepo.Create(ctx, item)
	if err != nil {
		return nil, apierror.ErrInternalServer
	}
	return item, nil
}

// GetVaultItems 检索用户的所有保险库项目。
func (s *VaultService) GetVaultItems(ctx context.Context, userID uuid.UUID) ([]core.VaultItem, error) {
	return s.vaultRepo.FindByUser(ctx, userID)
}

// GetVaultItemByID 通过其 ID 检索单个保险库项目，确保它属于该用户。
func (s *VaultService) GetVaultItemByID(ctx context.Context, id, userID uuid.UUID) (*core.VaultItem, error) {
	item, err := s.vaultRepo.FindByID(ctx, id)
	if err != nil {
		return nil, apierror.ErrNotFound
	}
	// 确保该项目属于请求用户
	if item.UserID != userID {
		return nil, apierror.ErrForbidden
	}
	return item, nil
}

// UpdateVaultItem 更新现有的保险库项目。
func (s *VaultService) UpdateVaultItem(ctx context.Context, item *core.VaultItem, userID uuid.UUID) (*core.VaultItem, error) {
	// 首先，验证该项目是否存在并属于该用户
	existingItem, err := s.GetVaultItemByID(ctx, item.ID, userID)
	if err != nil {
		return nil, err
	}

	// 确保用户 ID 不被更改
	item.UserID = existingItem.UserID

	// 如果更新请求中未提供类别，则保留现有类别。
	if item.Category == "" {
		item.Category = existingItem.Category
	}

	// 保留原始创建时间戳并更新修改时间戳。
	item.CreatedAt = existingItem.CreatedAt
	item.UpdatedAt = time.Now()

	err = s.vaultRepo.Update(ctx, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

// DeleteVaultItem 删除一个保险库项目。
func (s *VaultService) DeleteVaultItem(ctx context.Context, id, userID uuid.UUID) error {
	// 首先，验证该项目是否存在并属于该用户
	_, err := s.GetVaultItemByID(ctx, id, userID)
	if err != nil {
		return err
	}
	return s.vaultRepo.Delete(ctx, id)
}