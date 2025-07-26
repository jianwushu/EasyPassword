package core

import (
	"time"

	"github.com/google/uuid"
)

// User 表示系统中的一个用户。
type User struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Username   string    `gorm:"type:varchar(255);unique_index;not null"`
	Email      string    `gorm:"type:varchar(255);unique_index;not null"`
	AuthHash   string    `gorm:"type:text;not null"`
	MasterSalt []byte `gorm:"type:bytea;not null"`
	// for password reset
	ResetPasswordToken          *string
	ResetPasswordTokenExpiresAt *time.Time
	CreatedAt                   time.Time `gorm:"autoCreateTime"`
	UpdatedAt                   time.Time `gorm:"autoUpdateTime"`
}