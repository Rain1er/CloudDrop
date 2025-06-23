package routes

import (
	"clouddrop/config"
	"clouddrop/internal/handler"
	middleware2 "clouddrop/pkg/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(cfg *config.Config, db *gorm.DB) *gin.Engine {
	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// 中间件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware2.CORSMiddleware()) // 跨域中间件

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
		})
	})

	// 初始化处理器
	authHandler := handler.NewAuthHandler(cfg, db)
	webshellHandler := handler.NewWebShellHandler(cfg, db)

	// API路由组
	api := router.Group("/api/v1")
	{
		// 认证相关路由组
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
		}

		// webshell管理
		webshells := api.Group("/webshells")
		//webshells.Use(middleware2.AuthRequired(cfg.JWT.Secret)) // 鉴权
		{
			webshells.GET("/list", webshellHandler.List)

			webshells.GET("/:id", webshellHandler.Get)
			webshells.PUT("/:id", webshellHandler.Update)
			webshells.DELETE("/:id", webshellHandler.Delete)
			webshells.POST("/create", webshellHandler.Create)

			webshells.GET("/test/:id", webshellHandler.Test)
			webshells.POST("/batch-test", webshellHandler.BatchTest)

			webshells.GET("/BaseInfo/:id", webshellHandler.BaseInfo)
			webshells.POST("/ExecCommand/:id", webshellHandler.ExecCommand)
			webshells.POST("/ExecCode/:id", webshellHandler.ExecCode)
			webshells.POST("/ExecSql/:id", webshellHandler.ExecSql)

			webshells.POST("/FileZip/:id", webshellHandler.FileZip)
			webshells.POST("/FileUnZip/:id", webshellHandler.FileUnZip)

			webshells.POST("/FileList/:id", webshellHandler.FileList)
			webshells.POST("/FileShow/:id", webshellHandler.FileShow)

		}
	}

	return router
}
