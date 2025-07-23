package boltdb

import (
	"context"
	"easy-password-backend/internal/core"
	"encoding/json"

	"github.com/google/uuid"
	"go.etcd.io/bbolt"
)

var (
	userBucket     = []byte("users")
	vaultBucket    = []byte("vaults")
	usernameBucket = []byte("usernames")
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

// --- 用户存储库实现 ---

type userRepository struct {
	db *bbolt.DB
}

func (r *userRepository) Create(ctx context.Context, user *core.User) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		users := tx.Bucket(userBucket)
		usernames := tx.Bucket(usernameBucket)

		if usernames.Get([]byte(user.Username)) != nil {
			return &core.DuplicateEntryError{Field: "username"}
		}

		user.ID = uuid.New()
		encoded, err := json.Marshal(user)
		if err != nil {
			return err
		}

		if err := users.Put(user.ID[:], encoded); err != nil {
			return err
		}
		return usernames.Put([]byte(user.Username), user.ID[:])
	})
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*core.User, error) {
	var user core.User
	err := r.db.View(func(tx *bbolt.Tx) error {
		userID := tx.Bucket(usernameBucket).Get([]byte(username))
		if userID == nil {
			return core.ErrUserNotFound
		}

		userBytes := tx.Bucket(userBucket).Get(userID)
		if userBytes == nil {
			return core.ErrUserNotFound
		}
		return json.Unmarshal(userBytes, &user)
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

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