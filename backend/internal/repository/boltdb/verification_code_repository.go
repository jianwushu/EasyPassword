package boltdb

import (
	"context"
	"easy-password-backend/internal/core"
	"encoding/json"

	"go.etcd.io/bbolt"
)

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