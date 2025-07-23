package v1

import (
	"easy-password-backend/internal/apierror"
	"easy-password-backend/internal/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthHandler 处理与身份验证相关的 API 请求。
type AuthHandler struct {
	authService *auth.AuthService
}

// NewAuthHandler 创建一个新的 AuthHandler。
func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// RegisterRoutes 注册身份验证路由。
func (h *AuthHandler) RegisterRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1/auth")
	{
		v1.POST("/register", h.register)
		v1.POST("/login", h.login)
		v1.GET("/salt/:username", h.getSalt)
	}
}

type registerRequest struct {
	Username      string `json:"username" binding:"required,min=1"`
	MasterKeyHash string `json:"master_key_hash" binding:"required"`
	MasterSalt    string `json:"master_salt" binding:"required"`
}

type loginRequest struct {
	Username      string `json:"username" binding:"required"`
	MasterKeyHash string `json:"master_key_hash" binding:"required"`
}

type loginResponse struct {
	Token      string `json:"token"`
	MasterSalt string `json:"master_salt"`
}

func (h *AuthHandler) register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apierror.ErrInvalidRequest)
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req.Username, req.MasterKeyHash, req.MasterSalt)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user_id": user.ID,
	})
}

func (h *AuthHandler) login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apierror.ErrInvalidRequest)
		return
	}

	token, masterSalt, err := h.authService.Login(c.Request.Context(), req.Username, req.MasterKeyHash)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, loginResponse{
		Token:      token,
		MasterSalt: masterSalt,
	})
}

func (h *AuthHandler) getSalt(c *gin.Context) {
	username := c.Param("username")
	masterSalt, err := h.authService.GetMasterSalt(c.Request.Context(), username)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"master_salt": masterSalt})
}
