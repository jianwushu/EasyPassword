package apierror

import "net/http"

// APIError 表示用于 API 响应的结构化错误。
type APIError struct {
	Code    int    `json:"-"` // HTTP 状态码，在 JSON 响应体中忽略
	Message string `json:"message"`
}

// Error 使 APIError 满足错误接口。
func (e *APIError) Error() string {
	return e.Message
}

// New 创建一个新的 APIError。
func New(code int, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}

// 预定义的、可重用的错误实例。
var (
	ErrInvalidRequest     = New(http.StatusBadRequest, "Invalid request body")
	ErrUnauthorized       = New(http.StatusUnauthorized, "Authorization is required")
	ErrInvalidCredentials = New(http.StatusUnauthorized, "Invalid username or password")
	ErrInvalidToken       = New(http.StatusUnauthorized, "Invalid or expired token")
	ErrForbidden          = New(http.StatusForbidden, "Access denied")
	ErrNotFound           = New(http.StatusNotFound, "Resource not found")
	ErrUsernameExists     = New(http.StatusConflict, "Username already exists")
	ErrEmailExists             = New(http.StatusConflict, "Email already exists")
	ErrUserOrEmailExists       = New(http.StatusConflict, "Username or email already exists")
	ErrInvalidVerificationCode = New(http.StatusBadRequest, "Invalid verification code")
	ErrVerificationCodeExpired = New(http.StatusBadRequest, "Verification code has expired")
	ErrInternalServer          = New(http.StatusInternalServerError, "An unexpected error occurred")
)