package handler

import (
	"clouddrop/config"
	"clouddrop/internal/model"
	"clouddrop/internal/service"
	"time"

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

func (h *WebShellHandler) Identify(shellType string) service.Shell {
	var shell service.Shell // 声明接口类型的变量

	switch shellType {
	case "php":
		shell = &service.PHPShell{}
	case "java":
		shell = &service.JavaShell{}
	case "c#":
		shell = &service.CSharpShell{}
	case "asp":
		shell = &service.AspShell{}
	}

	return shell
}

// List 获取webshell列表
func (h *WebShellHandler) List(c *gin.Context) {

}

// Create 创建webshell
func (h *WebShellHandler) Create(c *gin.Context) {
	// 解析请求体到结构体，避免了单独获取每个参数，而且可以做数据验证，这里比Java的注解要容易理解，优雅
	var req CreateWebShellRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 创建WebShell记录，这里为了简单起见，直接使用service.PHPShell结构体。不去调用构造方法
	webshell := model.Web_shells{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      req.Name,
		URL:       req.URL,
		Password:  req.Password,
		Type:      req.Type,
		Encode:    req.Encode,
		Note:      req.Note,
	}

	// 插入数据到数据库
	result := h.db.Table("web_shells").Create(&webshell)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to create WebShell", "message": result.Error.Error()})
		return
	}

	// 返回成功信息
	c.JSON(201, gin.H{"message": "WebShell created successfully", "data": webshell})
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

// GetCurrentDirectory 获取当前目录
func (h *WebShellHandler) GetCurrentDirectory(c *gin.Context) {
	var currentDir string
	var err error
	// 从前端请求的ID查询数据库，获取WebShell的URL和密码
	id := c.Param("id")
	var webshell model.Web_shells
	if result := h.db.Where("id = ?", id).First(&webshell); result.Error != nil {
		c.JSON(404, gin.H{"error": "WebShell not found"})
		return
	}

	// 使用接口的多态特性，调用服务层获取当前目录
	shell := h.Identify(webshell.Type)
	currentDir, err = shell.GetCurrentDirectory(webshell.URL, webshell.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get current directory", "message": err.Error()})
		return
	}
	// 返回当前目录
	c.JSON(200, gin.H{"current_directory": currentDir})
}

func (h *WebShellHandler) ListFiles(c *gin.Context) {
	id := c.Param("id")
	var webshell model.Web_shells
	if res := h.db.Where("id = ?", id).First(&webshell); res.Error != nil {
		c.JSON(404, gin.H{"error": "WebShell not found"})
		return
	}
	shell := h.Identify(webshell.Type)
	listFiles, err := shell.ListFiles(webshell.URL, webshell.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to all files in the current directory", "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"current_directory_files": listFiles})
}

// ExecCommand 执行客户端发送的命令
func (h *WebShellHandler) ExecCommand(c *gin.Context) {
	id := c.Param("id")
	command := c.PostForm("command")

	var webshell model.Web_shells
	if res := h.db.Where("id = ?", id).First(&webshell); res.Error != nil {
		c.JSON(404, gin.H{"error": "WebShell not found"})
		return
	}
	shell := h.Identify(webshell.Type)
	info, err := shell.ExecCommand(webshell.URL, webshell.Password, command)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to ExecCommand", "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"command info": info})
}
