package boltdb

import (
	"context"
	"easy-password-backend/internal/core"
	"encoding/json"

	"github.com/google/uuid"
	"go.etcd.io/bbolt"
)

// --- 保险库存储库实现 ---

type vaultRepository struct {
	db *bbolt.DB
}

func (r *vaultRepository) Create(ctx context.Context, item *core.VaultItem) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		vaults := tx.Bucket(vaultBucket)
		item.ID = uuid.New()
		encoded, err := json.Marshal(item)
		if err != nil {
			return err
		}
		return vaults.Put(item.ID[:], encoded)
	})
}

func (r *vaultRepository) FindByID(ctx context.Context, id uuid.UUID) (*core.VaultItem, error) {
	var item core.VaultItem
	err := r.db.View(func(tx *bbolt.Tx) error {
		itemBytes := tx.Bucket(vaultBucket).Get(id[:])
		if itemBytes == nil {
			return core.ErrVaultItemNotFound
		}
		return json.Unmarshal(itemBytes, &item)
	})
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *vaultRepository) FindByUser(ctx context.Context, userID uuid.UUID) ([]core.VaultItem, error) {
	var items []core.VaultItem
	err := r.db.View(func(tx *bbolt.Tx) error {
		c := tx.Bucket(vaultBucket).Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var item core.VaultItem
			if err := json.Unmarshal(v, &item); err == nil {
				if item.UserID == userID {
					items = append(items, item)
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *vaultRepository) Update(ctx context.Context, item *core.VaultItem) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		vaults := tx.Bucket(vaultBucket)
		if existing := vaults.Get(item.ID[:]); existing == nil {
			return core.ErrVaultItemNotFound
		}
		encoded, err := json.Marshal(item)
		if err != nil {
			return err
		}
		return vaults.Put(item.ID[:], encoded)
	})
}

func (r *vaultRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		vaults := tx.Bucket(vaultBucket)
		if existing := vaults.Get(id[:]); existing == nil {
			return core.ErrVaultItemNotFound
		}
		return vaults.Delete(id[:])
	})
}
