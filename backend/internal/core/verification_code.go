package core

import "time"

// VerificationCode 用于存储发送给用户的邮箱验证码。
type VerificationCode struct {
	Email     string    `gorm:"type:varchar(255);primary_key"`
	Code      string    `gorm:"type:varchar(10);not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}