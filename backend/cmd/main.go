package main

import (
	config2 "clouddrop/config"
	"clouddrop/routes"
	"log"
)

func main() {
	// 1. 初始化配置
	cfg := config2.New()

	// 2. 初始化数据库实例，仅在第一次运行时会创建数据库文件
	db, err := config2.Initialize(cfg.Database.DSN)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// 3. 传入配置与数据库实例，进行初始化路由
	router := routes.SetupRouter(cfg, db)

	// 4. 启动服务器
	log.Printf("CloudDrop server starting on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
