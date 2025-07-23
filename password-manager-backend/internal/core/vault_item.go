package core

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// VaultItem 表示用户保险库中的一个加密项目。
type VaultItem struct {
	ID             uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID         uuid.UUID       `gorm:"type:uuid;not null"`
	EncryptedData  json.RawMessage `gorm:"type:jsonb;not null"`
	Category       string          `gorm:"type:varchar(100);index"` // 添加类别字段
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}