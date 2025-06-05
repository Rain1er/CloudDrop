package handlers

import (
	"clouddrop/internal/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WebShellHandler struct {
	db  *gorm.DB       // 数据库连接
	cfg *config.Config // 配置
}

// AuthHandler 创建认证处理器
func NewWebShellHandler(cfg *config.Config, db *gorm.DB) *WebShellHandler {
	return &WebShellHandler{
		cfg: cfg,
		db:  db,
	}
}

// 定义新的 CreateWebShellRequest 数据结构
type CreateWebShellRequest struct {
	Name     string `json:"name" binding:"required"`
	URL      string `json:"url" binding:"required"`
	Password string `json:"password" binding:"required"`
	Type     string `json:"type"`
	Encode   string `json:"encode"`
	Note     string `json:"note"`
}

// List 获取webshell列表
func (h *WebShellHandler) List(c *gin.Context) {

}

// Create 创建webshell
func (h *WebShellHandler) Create(c *gin.Context) {

}

// Get 获取单个webshell
func (h *WebShellHandler) Get(c *gin.Context) {

}

// Update 更新webshell
func (h *WebShellHandler) Update(c *gin.Context) {

}

// Delete 删除webshell
func (h *WebShellHandler) Delete(c *gin.Context) {

}

// Test 测试单个WebShell有效性
func (h *WebShellHandler) Test(c *gin.Context) {
}

// BatchTest 批量测试WebShell连接
func (h *WebShellHandler) BatchTest(c *gin.Context) {

}
