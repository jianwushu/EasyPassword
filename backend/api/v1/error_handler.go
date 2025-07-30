package v1

import (
	"easy-password-backend/internal/apierror"
	"errors"

	"github.com/gin-gonic/gin"
)

// handleError 集中处理所有 API 处理程序的错误。
func handleError(c *gin.Context, err error) {
	var apiErr *apierror.APIError
	if errors.As(err, &apiErr) {
		c.JSON(apiErr.Code, gin.H{"error": apiErr.Message})
		return
	}

	// 对于任何其他错误，返回一个通用的 500 内部服务器错误。
	c.JSON(apierror.ErrInternalServer.Code, gin.H{"error": apierror.ErrInternalServer.Message})
}