package handler

import (
	"clouddrop/config"
	"clouddrop/internal/model"
	"net/http"
	"time"

	"crypto/md5"
	"encoding/hex"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthHandler struct {
	cfg *config.Config
	db  *gorm.DB
}

// AuthHandler 创建认证处理器
func NewAuthHandler(cfg *config.Config, db *gorm.DB) *AuthHandler {
	return &AuthHandler{
		cfg: cfg,
		db:  db,
	}
}

type LoginRequest struct {
	// binding:"required" 是 Go 的 Gin 框架 中用于验证数据的标签，类似于Java中的注解
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Role     string `json:"role"`
	} `json:"user"`
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查询用户
	var user model.Users
	if result := h.db.Where("username = ?", req.Username).First(&user); result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	hash := md5.Sum([]byte(req.Password))       // 计算 MD5 哈希
	passwordHash := hex.EncodeToString(hash[:]) // 转换为十六进制字符串
	// 验证密码
	if user.Password != passwordHash {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLogin = &now
	h.db.Save(&user)

	// 生成JWT令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * time.Duration(h.cfg.JWT.Expire)).Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.cfg.JWT.Secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	// 构建响应，其中包含JWT令牌和用户信息
	resp := LoginResponse{
		Token: tokenString,
	}
	resp.User.ID = int(user.ID)
	resp.User.Username = user.Username
	resp.User.Role = user.Role

	c.JSON(http.StatusOK, resp)
}
