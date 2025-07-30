package service

import (
	"context"
	"easy-password-backend/internal/apierror"
	"easy-password-backend/internal/core"
	"log/slog"
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
	slog.Info("Creating new vault item", "user_id", item.UserID)
	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now
	err := s.vaultRepo.Create(ctx, item)
	if err != nil {
		slog.Error("Failed to create vault item", "user_id", item.UserID, "error", err)
		return nil, apierror.ErrInternalServer
	}
	slog.Info("Vault item created successfully", "item_id", item.ID, "user_id", item.UserID)
	return item, nil
}

// GetVaultItems 检索用户的所有保险库项目。
func (s *VaultService) GetVaultItems(ctx context.Context, userID uuid.UUID) ([]core.VaultItem, error) {
	slog.Info("Fetching vault items for user", "user_id", userID)
	items, err := s.vaultRepo.FindByUser(ctx, userID)
	if err != nil {
		slog.Error("Failed to fetch vault items", "user_id", userID, "error", err)
		// 在这种情况下，我们可能不想向客户端暴露内部服务器错误，
		// 但出于日志目的，记录它是很重要的。
		// 返回一个空切片和 nil 错误，或者一个特定的应用错误。
		// 为了简单起见，我们现在只记录并返回。
	}
	slog.Info("Fetched vault items", "user_id", userID, "count", len(items))
	return items, err
}

// GetVaultItemByID 通过其 ID 检索单个保险库项目，确保它属于该用户。
func (s *VaultService) GetVaultItemByID(ctx context.Context, id, userID uuid.UUID) (*core.VaultItem, error) {
	slog.Info("Fetching vault item by ID", "item_id", id, "user_id", userID)
	item, err := s.vaultRepo.FindByID(ctx, id)
	if err != nil {
		slog.Warn("Vault item not found by ID", "item_id", id, "error", err)
		return nil, apierror.ErrNotFound
	}
	// 确保该项目属于请求用户
	if item.UserID != userID {
		slog.Warn("User forbidden to access vault item", "item_id", id, "user_id", userID, "owner_id", item.UserID)
		return nil, apierror.ErrForbidden
	}
	slog.Info("Vault item fetched successfully", "item_id", id)
	return item, nil
}

// UpdateVaultItem 更新现有的保险库项目。
func (s *VaultService) UpdateVaultItem(ctx context.Context, item *core.VaultItem, userID uuid.UUID) (*core.VaultItem, error) {
	slog.Info("Updating vault item", "item_id", item.ID, "user_id", userID)
	// 首先，验证该项目是否存在并属于该用户
	existingItem, err := s.GetVaultItemByID(ctx, item.ID, userID)
	if err != nil {
		// GetVaultItemByID 已经记录了错误
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
		slog.Error("Failed to update vault item", "item_id", item.ID, "error", err)
		return nil, err
	}

	slog.Info("Vault item updated successfully", "item_id", item.ID)
	return item, nil
}

// DeleteVaultItem 删除一个保险库项目。
func (s *VaultService) DeleteVaultItem(ctx context.Context, id, userID uuid.UUID) error {
	slog.Info("Deleting vault item", "item_id", id, "user_id", userID)
	// 首先，验证该项目是否存在并属于该用户
	_, err := s.GetVaultItemByID(ctx, id, userID)
	if err != nil {
		// GetVaultItemByID 已经记录了错误
		return err
	}
	err = s.vaultRepo.Delete(ctx, id)
	if err != nil {
		slog.Error("Failed to delete vault item", "item_id", id, "error", err)
		return err
	}
	slog.Info("Vault item deleted successfully", "item_id", id)
	return nil
}