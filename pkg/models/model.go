package models

import (
	"time"
)

// Model 表示AI模型的数据结构
type Model struct {
	ModelID    string    `json:"model_id" gorm:"primaryKey;type:varchar(64)"`
	Name       string    `json:"name" gorm:"type:varchar(255);not null;uniqueIndex" binding:"required"`
	Endpoint   string    `json:"endpoint" gorm:"type:varchar(255);not null" binding:"required"`
	APIKey     string    `json:"api_key" gorm:"type:varchar(255);not null" binding:"required"`
	Timeout    int       `json:"timeout" gorm:"not null" binding:"required"`
	Type       string    `json:"type" gorm:"type:varchar(255);not null" binding:"required"`
	Dimensions int       `json:"dimensions" gorm:"not null" binding:"required"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// CreateModelRequest 创建模型的请求结构
type CreateModelRequest struct {
	Name       string `json:"name" binding:"required"`
	Endpoint   string `json:"endpoint" binding:"required"`
	APIKey     string `json:"api_key" binding:"required"`
	Timeout    int    `json:"timeout" binding:"required"`
	Type       string `json:"type" binding:"required"`
	Dimensions int    `json:"dimensions" binding:"required"`
}

// UpdateModelRequest 更新模型的请求结构
type UpdateModelRequest struct {
	Name       *string `json:"name"`
	Endpoint   *string `json:"endpoint"`
	APIKey     *string `json:"api_key"`
	Timeout    *int    `json:"timeout"`
	Type       *string `json:"type"`
	Dimensions *int    `json:"dimensions"`
}

// ChatMessage OpenAI风格的对话消息结构体
type ChatMessage struct {
	Role    string `json:"role" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// ChatRequest OpenAI风格的对话请求结构体
type ChatRequest struct {
	Model    string        `json:"model" binding:"required"`
	Messages []ChatMessage `json:"messages" binding:"required"`
}

// TableName 指定表名
func (Model) TableName() string {
	return "t_model"
}
