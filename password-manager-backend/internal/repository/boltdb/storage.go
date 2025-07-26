package boltdb

import (
	"context"
	"easy-password-backend/internal/core"
	"encoding/json"

	"github.com/google/uuid"
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

// --- 用户存储库实现 ---

type userRepository struct {
	db *bbolt.DB
}

func (r *userRepository) Create(ctx context.Context, user *core.User) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		users := tx.Bucket(userBucket)
		usernames := tx.Bucket(usernameBucket)
		emails := tx.Bucket(emailBucket)

		if usernames.Get([]byte(user.Username)) != nil {
			return &core.DuplicateEntryError{Field: "username"}
		}
		if emails.Get([]byte(user.Email)) != nil {
			return &core.DuplicateEntryError{Field: "email"}
		}

		user.ID = uuid.New()
		encoded, err := json.Marshal(user)
		if err != nil {
			return err
		}

		if err := users.Put(user.ID[:], encoded); err != nil {
			return err
		}
		if err := usernames.Put([]byte(user.Username), user.ID[:]); err != nil {
			return err
		}
		return emails.Put([]byte(user.Email), user.ID[:])
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

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*core.User, error) {
	var user core.User
	err := r.db.View(func(tx *bbolt.Tx) error {
		userID := tx.Bucket(emailBucket).Get([]byte(email))
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

func (r *userRepository) FindByResetPasswordToken(ctx context.Context, token string) (*core.User, error) {
	var foundUser *core.User
	err := r.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(userBucket)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var user core.User
			if err := json.Unmarshal(v, &user); err != nil {
				// 忽略无法解析的条目，或者记录日志
				continue
			}
			if user.ResetPasswordToken != nil && *user.ResetPasswordToken == token {
				foundUser = &user
				return nil // 找到后停止遍历
			}
		}
		// 如果遍历完都没有找到
		return core.ErrUserNotFound
	})

	if err != nil {
		return nil, err
	}
	return foundUser, nil
}

func (r *userRepository) Update(ctx context.Context, user *core.User) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		users := tx.Bucket(userBucket)

		// 确保用户存在
		if existing := users.Get(user.ID[:]); existing == nil {
			return core.ErrUserNotFound
		}

		encoded, err := json.Marshal(user)
		if err != nil {
			return err
		}
		return users.Put(user.ID[:], encoded)
	})
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

// --- 验证码存储库实现 ---

type verificationCodeRepository struct {
	db *bbolt.DB
}

func (r *verificationCodeRepository) Create(ctx context.Context, vc *core.VerificationCode) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(verificationCodeBucket)
		encoded, err := json.Marshal(vc)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(vc.Email), encoded)
	})
}

func (r *verificationCodeRepository) Find(ctx context.Context, email string) (*core.VerificationCode, error) {
	var vc core.VerificationCode
	err := r.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(verificationCodeBucket)
		vcBytes := bucket.Get([]byte(email))
		if vcBytes == nil {
			return core.ErrVerificationCodeNotFound
		}
		return json.Unmarshal(vcBytes, &vc)
	})
	if err != nil {
		return nil, err
	}
	return &vc, nil
}

func (r *verificationCodeRepository) Delete(ctx context.Context, email string) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(verificationCodeBucket)
		return bucket.Delete([]byte(email))
	})
}