package core

import "errors"

// 存储库的预定义错误
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrVaultItemNotFound = errors.New("vault item not found")
	ErrVerificationCodeNotFound = errors.New("verification code not found")
)

// 当违反唯一约束时返回 DuplicateEntryError。
type DuplicateEntryError struct {
	Field string
}

func (e *DuplicateEntryError) Error() string {
	return "duplicate entry for field: " + e.Field
}