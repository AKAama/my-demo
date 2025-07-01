package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Model 表示AI模型的数据结构
type Model struct {
	ModelID   string    `json:"model_id" gorm:"primaryKey;type:varchar(64)"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null;uniqueIndex" binding:"required"`
	Endpoint  string    `json:"endpoint" gorm:"type:varchar(255);not null" binding:"required"`
	APIKey    string    `json:"api_key" gorm:"type:varchar(255);not null" binding:"required"`
	Timeout   int       `json:"timeout" gorm:"not null" binding:"required"`
	Type      string    `json:"type" gorm:"type:varchar(255);not null" binding:"required"`
	Dimension int       `json:"dimension" gorm:"not null" binding:"required"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// CreateModelRequest 创建模型的请求结构
type CreateModelRequest struct {
	Name      string `json:"name" binding:"required"`
	Endpoint  string `json:"endpoint" binding:"required"`
	APIKey    string `json:"api_key" binding:"required"`
	Timeout   int    `json:"timeout" binding:"required"`
	Type      string `json:"type" binding:"required"`
	Dimension int    `json:"dimension" binding:"required"`
}

// UpdateModelRequest 更新模型的请求结构
type UpdateModelRequest struct {
	Name      *string `json:"name"`
	Endpoint  *string `json:"endpoint"`
	APIKey    *string `json:"api_key"`
	Timeout   *int    `json:"timeout"`
	Type      *string `json:"type"`
	Dimension *int    `json:"dimension"`
}

// BeforeCreate GORM钩子，在创建前生成UUID
func (m *Model) BeforeCreate(*gorm.DB) error {
	if m.ModelID == "" {
		m.ModelID = uuid.New().String()
	}
	return nil
}

// TableName 指定表名
func (Model) TableName() string {
	return "models"
}
