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
		v1.POST("/salt", h.getSalt)
		v1.POST("/send-verification-code", h.sendVerificationCode)
		v1.POST("/request-password-reset", h.requestPasswordReset)
		v1.POST("/reset-password", h.resetPassword)
	}
}

type registerRequest struct {
	Username      string `json:"username" binding:"required,min=1"`
	Email         string `json:"email" binding:"required,email"`
	MasterKeyHash string `json:"master_key_hash" binding:"required"`
	MasterSalt    string `json:"master_salt" binding:"required"`
	Code          string `json:"code" binding:"required,len=6"`
}

type loginRequest struct {
	Identifier    string `json:"identifier" binding:"required"`
	MasterKeyHash string `json:"master_key_hash" binding:"required"`
}

type loginResponse struct {
	Username   string `json:"username"`
	Token      string `json:"token"`
	MasterSalt string `json:"master_salt"`
}

type getSaltRequest struct {
	Identifier string `json:"identifier" binding:"required"`
}

type sendVerificationCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type requestPasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type resetPasswordRequest struct {
	Token            string `json:"token" binding:"required"`
	NewMasterKeyHash string `json:"new_master_key_hash" binding:"required"`
	NewMasterSalt    string `json:"new_master_salt" binding:"required"`
}

func (h *AuthHandler) register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apierror.ErrInvalidRequest)
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req.Username, req.Email, req.MasterKeyHash, req.MasterSalt, req.Code)
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

	token, masterSalt, username, err := h.authService.Login(c.Request.Context(), req.Identifier, req.MasterKeyHash)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, loginResponse{
		Username:   username,
		Token:      token,
		MasterSalt: masterSalt,
	})
}

func (h *AuthHandler) getSalt(c *gin.Context) {
	var req getSaltRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apierror.ErrInvalidRequest)
		return
	}
	masterSalt, err := h.authService.GetMasterSalt(c.Request.Context(), req.Identifier)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"master_salt": masterSalt})
}

func (h *AuthHandler) sendVerificationCode(c *gin.Context) {
	var req sendVerificationCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apierror.ErrInvalidRequest)
		return
	}

	err := h.authService.SendVerificationCode(c.Request.Context(), req.Email)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification code sent successfully"})
}

func (h *AuthHandler) requestPasswordReset(c *gin.Context) {
	var req requestPasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apierror.ErrInvalidRequest)
		return
	}

	err := h.authService.RequestPasswordReset(c.Request.Context(), req.Email)
	if err != nil {
		handleError(c, err)
		return
	}

	// 出于安全考虑，即使找不到电子邮件，也始终返回成功的响应，以防止用户枚举攻击。
	c.JSON(http.StatusOK, gin.H{"message": "If an account with that email exists, a password reset link has been sent."})
}

func (h *AuthHandler) resetPassword(c *gin.Context) {
	var req resetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apierror.ErrInvalidRequest)
		return
	}

	err := h.authService.ResetPassword(c.Request.Context(), req.Token, req.NewMasterKeyHash, req.NewMasterSalt)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password has been reset successfully."})
}
