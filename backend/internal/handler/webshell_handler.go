package handler

import (
	"clouddrop/config"
	"clouddrop/internal/model"
	"clouddrop/internal/service"
	"strconv"
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

// 定义 WebShellRequest 数据结构
type WebShellRequest struct {
	Name     string `json:"name" binding:"required"`
	URL      string `json:"url" binding:"required"`
	Password string `json:"password" binding:"required"`
	Type     string `json:"type" binding:"required"`
	Encode   string `json:"encode" binding:"required"`
	Note     string `json:"note"`
}

func (h *WebShellHandler) GetType(shellType string) service.Shell {
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
	var webshells []model.Web_shells
	if result := h.db.Find(&webshells); result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve WebShells", "message": result.Error.Error()})
		return
	}
	c.JSON(200, gin.H{"results": webshells})

}

// Create 创建webshell
func (h *WebShellHandler) Create(c *gin.Context) {
	// 解析请求体到结构体，避免了单独获取每个参数，而且可以做数据验证，这里比Java的注解要容易理解，优雅
	var req WebShellRequest
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
	var webshell model.Web_shells
	id := c.Param("id")
	if result := h.db.Where("id = ?", id).First(&webshell); result.Error != nil {
		c.JSON(404, gin.H{"error": "Webshell Not found"})
		return
	}

	c.JSON(200, gin.H{"data": webshell})

}

// Update 更新webshell
func (h *WebShellHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req WebShellRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 在更新之前先要查询该shell是否存在
	var webshell model.Web_shells
	if result := h.db.Where("id = ?", id).First(&webshell); result.Error != nil {
		c.JSON(404, gin.H{"error": "WebShell not found"})
		return
	}

	// Update only specific fields, keeping CreatedAt unchanged
	updates := map[string]interface{}{
		"updated_at": time.Now(),
		"name":       req.Name,
		"url":        req.URL,
		"password":   req.Password,
		"type":       req.Type,
		"encode":     req.Encode,
		"note":       req.Note,
	}

	if result := h.db.Model(&webshell).Updates(updates); result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to update WebShell", "message": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "WebShell updated successfully", "data": webshell})
}

// Delete 删除webshell
func (h *WebShellHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	// First check if the WebShell exists
	var webshell model.Web_shells
	if result := h.db.Where("id = ?", id).First(&webshell); result.Error != nil {
		c.JSON(404, gin.H{"error": "WebShell not found"})
		return
	}

	// Delete the WebShell
	if result := h.db.Delete(&model.Web_shells{}, id); result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to delete WebShell", "message": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "WebShell deleted successfully"})
}

// Test 测试单个WebShell有效性
// 对于打开的webshell，每15分钟刷新一次session，避免session过期。可以使用前端Ajax请求实现
// 具体来说，第一次打开会话获取一次session，如果用户停留页面，那么每15分钟通过ajax请求刷新一次session
func (h *WebShellHandler) Test(c *gin.Context) {
	// 将webshell id与session进行绑定，放到全局变量中
	// 1. 请求一次webshell，除了获取session之外什么也不干
	id := c.Param("id")
	intID, _ := strconv.Atoi(id)

	var webshell model.Web_shells
	if result := h.db.Where("id = ?", id).First(&webshell); result.Error != nil {
		c.JSON(404, gin.H{"error": "WebShell not found"})
		return
	}
	shellHandler := h.GetType(webshell.Type)
	res, err := shellHandler.FreshSession(intID, webshell.URL, webshell.Password)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to Test", "message": err.Error()})
	}

	c.JSON(200, gin.H{"results": res})
}

// BatchTest 批量测试WebShell连接
func (h *WebShellHandler) BatchTest(c *gin.Context) {
	// Get all WebShells from database
	var webshells []model.Web_shells
	if result := h.db.Find(&webshells); result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve WebShells", "message": result.Error.Error()})
		return
	}

	// Test each WebShell and collect results
	results := make(map[int]bool)
	// for _, webshell := range webshells {
	// 	// Use the appropriate shell handler based on type
	// 	shellHandler := h.GetType(webshell.Type)
	// 	// Test connection and store result
	// 	_, err := shellHandler.BaseInfo(webshell.URL, webshell.Password)
	// 	results[webshell.ID] = err == nil
	// }

	c.JSON(200, gin.H{"results": results})
}

// BaseInfo 获取系统信息
func (h *WebShellHandler) BaseInfo(c *gin.Context) {
	// 拿到webshell信息
	var info string
	var err error
	var webshell model.Web_shells
	// 从前端请求的ID查询数据库，获取WebShell的URL和密码
	id := c.Param("id")
	intID, _ := strconv.Atoi(id)
	if result := h.db.Where("id = ?", id).First(&webshell); result.Error != nil {
		c.JSON(404, gin.H{"error": "WebShell not found"})
		return
	}

	// 使用接口的多态特性，调用服务层
	shellHandler := h.GetType(webshell.Type)
	info, err = shellHandler.BaseInfo(intID, webshell.URL, webshell.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get current directory", "message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"info": info})
}

// FileList 列出目录下的文件
func (h *WebShellHandler) FileList(c *gin.Context) {
	id := c.Param("id")
	intID, _ := strconv.Atoi(id)
	path := c.PostForm("path")
	var webshell model.Web_shells
	if res := h.db.Where("id = ?", id).First(&webshell); res.Error != nil {
		c.JSON(404, gin.H{"error": "WebShell not found"})
		return
	}
	shellHandler := h.GetType(webshell.Type)
	files, err := shellHandler.FileList(intID, path, webshell.URL, webshell.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to all files in the target directory", "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"target_directory_files": files})
}

// ExecCommand 执行客户端发送的命令
func (h *WebShellHandler) ExecCommand(c *gin.Context) {
	id := c.Param("id")
	intID, _ := strconv.Atoi(id)
	command := c.PostForm("command")

	var webshell model.Web_shells
	if res := h.db.Where("id = ?", id).First(&webshell); res.Error != nil {
		c.JSON(404, gin.H{"error": "WebShell not found"})
		return
	}
	shellHandler := h.GetType(webshell.Type)
	// Todo 单引号对于win可能会出错，需要在CMD.php中处理引号问题
	info, err := shellHandler.ExecCommand(intID, "'"+command+"'", webshell.URL, webshell.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to ExecCommand", "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"command info": info})
}
