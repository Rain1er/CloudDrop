package config

import (
	"clouddrop/internal/model"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Initialize 初始化数据库连接并自动迁移模型
func Initialize(dsn string) (*gorm.DB, error) {
	// 配置GORM日志
	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // 慢SQL阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略记录未找到错误
			Colorful:                  true,        // 彩色打印
		},
	)

	// 连接数据库
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}
	// 仅仅在程序第一次运行时创建数据库
	if _, err := os.Stat(dsn); os.IsNotExist(err) {
		log.Printf("数据库文件 %s 不存在，将创建新数据库", dsn)
		// 自动迁移数据表结构
		err = db.AutoMigrate(&model.Users{}, &model.Web_shells{})
		if err != nil {
			return nil, err
		}

		// 创建默认管理员用户
		var adminCount int64
		db.Model(&model.Users{}).Where("username = ?", "admin").Count(&adminCount)
		if adminCount == 0 {
			// 添加初始管理员账户
			adminUser := model.Users{
				Username: "admin",
				Password: "3cc87212bcd411686a3b9e547d47fc51", // raindrop
				Role:     "admin",
				IsActive: true,
			}
			db.Create(&adminUser)
			log.Println("默认管理员用户已创建")
		}

		log.Println("数据库初始化完成")
	}

	return db, nil
}
