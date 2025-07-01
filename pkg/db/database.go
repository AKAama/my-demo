package db

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"myapi/config"
	"myapi/pkg/models" // 新增导入
)

var gormDB *gorm.DB
var tidbOnce sync.Once

func InitTiDB(cfg *config.GlobalConfig) error {
	var err error
	tidbOnce.Do(func() {
		//zap.S().Infof(cfg.DSN())
		gormDB, err = gorm.Open(mysql.New(mysql.Config{
			DSN: cfg.DBConfig.DSN(),
		}), &gorm.Config{
			NowFunc: func() time.Time {
				ti, _ := time.LoadLocation("Asia/Shanghai")
				return time.Now().In(ti)
			},
			Logger: logger.Default.LogMode(logger.Silent),
		})

		if err != nil {
			return
		}

		if err != nil {
			return
		}

		if err = initTiDB(cfg.DBConfig.MaxConnections); err != nil {
			return
		}
		zap.S().Debug("database init finished...")
	})
	return err
}

func initTiDB(maxConnections int) error {
	maxValue := 100
	if maxConnections > 0 {
		maxValue = maxConnections
	}

	err := gormDB.Use(
		dbresolver.Register(dbresolver.Config{}).
			SetMaxOpenConns(maxValue),
	)
	if err != nil {
		return err
	}
	// 添加模型表的自动迁移
	return gormDB.AutoMigrate(&models.Model{})
}

func GetDB() *gorm.DB {
	return gormDB
}

func GetDBWithContext(ctx context.Context) *gorm.DB {
	return gormDB.WithContext(ctx)
}
