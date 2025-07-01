package server

import (
	"github.com/gin-gonic/gin"
)

func InitRouter(engine *gin.Engine) {
	// 创建模型处理器
	modelHandler := NewModelHandler()

	// API路由组
	api := engine.Group("/api/v1")
	{
		// 模型管理路由
		models := api.Group("/models")
		{
			models.POST("/create", modelHandler.CreateModel) // 创建模型
			models.GET("/get", modelHandler.GetModels)       // 获取模型列表
			models.GET("/:id", modelHandler.GetModel)        // 获取单个模型
			models.PUT("/:id", modelHandler.UpdateModel)     // 更新模型
			models.DELETE("/:id", modelHandler.DeleteModel)  // 删除模型
		}
	}
}
