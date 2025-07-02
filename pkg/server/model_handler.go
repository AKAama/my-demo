package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"myapi/pkg/db"
	"myapi/pkg/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ModelHandler 模型相关的处理器
type ModelHandler struct{}

// NewModelHandler 创建新的模型处理器
func NewModelHandler() *ModelHandler {
	return &ModelHandler{}
}

// CreateModel 创建模型
func (h *ModelHandler) CreateModel(c *gin.Context) {
	var req models.CreateModelRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		zap.S().Errorf("创建模型参数绑定错误: %v", err)
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(400, "参数错误: "+err.Error()))
		return
	}

	// 转换为Model结构
	model := models.Model{
		Name:      req.Name,
		Endpoint:  req.Endpoint,
		APIKey:    req.APIKey,
		Timeout:   req.Timeout,
		Type:      req.Type,
		Dimension: req.Dimension,
	}

	ctx := context.Background()
	database := db.GetDBWithContext(ctx)

	// 检查模型名是否已存在
	var existingModel models.Model
	if err := database.Where("name = ?", model.Name).First(&existingModel).Error; err == nil {
		zap.S().Warnf("尝试创建重复模型名称: %s", model.Name)
		c.JSON(http.StatusConflict, models.NewErrorResponse(409, "模型名称已存在"))
		return
	}

	// 创建模型
	if err := database.Create(&model).Error; err != nil {
		zap.S().Errorf("创建模型失败: %v", err)
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(500, "创建模型失败: "+err.Error()))
		return
	}

	zap.S().Infof("成功创建模型: %s, ID: %s", model.Name, model.ModelID)
	c.JSON(http.StatusOK, models.NewSuccessResponse(model, "成功创建模型"))
}

// GetModel 获取单个模型
func (h *ModelHandler) GetModel(c *gin.Context) {
	modelID := c.Param("id")

	ctx := context.Background()
	database := db.GetDBWithContext(ctx)

	var model models.Model
	if err := database.Where("model_id = ?", modelID).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zap.S().Warnf("请求的模型不存在: %s", modelID)
			c.JSON(http.StatusNotFound, models.NewErrorResponse(404, "模型不存在"))
		} else {
			zap.S().Errorf("查询模型失败: %v", err)
			c.JSON(http.StatusInternalServerError, models.NewErrorResponse(500, "查询模型失败: "+err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(model, "查询成功"))
}

// GetModels 获取模型列表
func (h *ModelHandler) GetModels(c *gin.Context) {
	ctx := context.Background()
	database := db.GetDBWithContext(ctx)

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	var modelList []models.Model
	var total int64

	// 查询总数
	if err := database.Model(&models.Model{}).Count(&total).Error; err != nil {
		zap.S().Errorf("查询模型总数失败: %v", err)
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(500, "查询模型总数失败: "+err.Error()))
		return
	}

	// 分页查询
	if err := database.Limit(pageSize).Offset((page - 1) * pageSize).Find(&modelList).Error; err != nil {
		zap.S().Errorf("查询模型列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(500, "查询模型列表失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(gin.H{
		"list":      modelList,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}, "查询模型列表成功"))
}

// UpdateModel 更新模型
func (h *ModelHandler) UpdateModel(c *gin.Context) {
	modelID := c.Param("id")

	ctx := context.Background()
	database := db.GetDBWithContext(ctx)

	var model models.Model
	if err := database.Where("model_id = ?", modelID).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zap.S().Warnf("尝试更新不存在的模型: %s", modelID)
			c.JSON(http.StatusNotFound, models.NewErrorResponse(404, "模型不存在"))
		} else {
			zap.S().Errorf("查询模型失败: %v", err)
			c.JSON(http.StatusInternalServerError, models.NewErrorResponse(500, "查询模型失败: "+err.Error()))
		}
		return
	}

	var req models.UpdateModelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.S().Errorf("更新模型参数绑定错误: %v", err)
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(400, "参数错误: "+err.Error()))
		return
	}

	// 检查名称是否与其他模型冲突
	if req.Name != nil && *req.Name != model.Name {
		var existingModel models.Model
		if err := database.Where("name = ? AND model_id != ?", *req.Name, modelID).First(&existingModel).Error; err == nil {
			zap.S().Warnf("尝试更新为已存在的模型名称: %s", *req.Name)
			c.JSON(http.StatusConflict, models.NewErrorResponse(409, "模型名称已存在"))
			return
		}
	}

	// 只更新提供的字段
	if req.Name != nil {
		model.Name = *req.Name
	}
	if req.Endpoint != nil {
		model.Endpoint = *req.Endpoint
	}
	if req.APIKey != nil {
		model.APIKey = *req.APIKey
	}
	if req.Timeout != nil {
		model.Timeout = *req.Timeout
	}
	if req.Type != nil {
		model.Type = *req.Type
	}
	if req.Dimension != nil {
		model.Dimension = *req.Dimension
	}

	if err := database.Save(&model).Error; err != nil {
		zap.S().Errorf("更新模型失败: %v", err)
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(500, "更新模型失败: "+err.Error()))
		return
	}

	zap.S().Infof("成功更新模型: %s, ID: %s", model.Name, model.ModelID)
	c.JSON(http.StatusOK, models.NewSuccessResponse(model, "成功更新模型"))
}

// DeleteModel 删除模型
func (h *ModelHandler) DeleteModel(c *gin.Context) {
	modelID := c.Param("id")

	ctx := context.Background()
	database := db.GetDBWithContext(ctx)

	var model models.Model
	if err := database.Where("model_id = ?", modelID).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zap.S().Warnf("尝试删除不存在的模型: %s", modelID)
			c.JSON(http.StatusNotFound, models.NewErrorResponse(404, "模型不存在"))
		} else {
			zap.S().Errorf("查询模型失败: %v", err)
			c.JSON(http.StatusInternalServerError, models.NewErrorResponse(500, "查询模型失败: "+err.Error()))
		}
		return
	}

	if err := database.Delete(&model).Error; err != nil {
		zap.S().Errorf("删除模型失败: %v", err)
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(500, "删除模型失败: "+err.Error()))
		return
	}

	zap.S().Infof("成功删除模型: %s, ID: %s", model.Name, model.ModelID)
	c.JSON(http.StatusOK, models.NewSuccessResponse(model, "成功删除模型"))
}

// ChatWithModel 大模型对话接口
func (h *ModelHandler) ChatWithModel(c *gin.Context) {
	modelID := c.Param("id")
	ctx := context.Background()
	database := db.GetDBWithContext(ctx)

	// 查找模型
	var model models.Model
	if err := database.Where("model_id = ?", modelID).First(&model).Error; err != nil {
		c.JSON(http.StatusNotFound, models.NewErrorResponse(404, "模型不存在"))
		return
	}

	// 解析提问内容
	var req models.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(400, "参数错误: "+err.Error()))
		return
	}

	// 构造大模型API请求
	payload, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(500, "请求序列化失败: "+err.Error()))
		return
	}
	httpReq, err := http.NewRequestWithContext(ctx, "POST", model.Endpoint, bytes.NewBuffer(payload))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(500, "请求创建失败: "+err.Error()))
		return
	}
	httpReq.Header.Set("Content-Type", "application/json")

	apiKey := strings.TrimSpace(model.APIKey)
	if apiKey == "" {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(500, "API Key 为空"))
		return
	}
	if strings.ContainsAny(apiKey, "\r\n\t ") {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(500, "API Key 包含非法字符"))
		return
	}
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	// 发起请求
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(500, "大模型请求失败: "+err.Error()))
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	body, _ := io.ReadAll(resp.Body)

	// 直接返回大模型响应
	c.Data(resp.StatusCode, "application/json", body)
}
