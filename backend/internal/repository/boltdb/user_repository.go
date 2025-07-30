package boltdb

import (
	"context"
	"easy-password-backend/internal/core"
	"encoding/json"

	"github.com/google/uuid"
	"go.etcd.io/bbolt"
)

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
