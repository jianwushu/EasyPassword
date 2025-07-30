package v1

import (
	"easy-password-backend/internal/apierror"
	"easy-password-backend/internal/core"
	"easy-password-backend/internal/service"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// VaultHandler 处理与保险库相关的 API 请求。
type VaultHandler struct {
	vaultService *service.VaultService
}

// NewVaultHandler 创建一个新的 VaultHandler。
func NewVaultHandler(vaultService *service.VaultService) *VaultHandler {
	return &VaultHandler{vaultService: vaultService}
}

// RegisterRoutes 注册保险库路由。
func (h *VaultHandler) RegisterRoutes(router *gin.RouterGroup) {
	vault := router.Group("/vault")
	{
		vault.POST("/items", h.createItem)
		vault.GET("/items", h.getItems)
		vault.PUT("/items/:id", h.updateItem)
		vault.DELETE("/items/:id", h.deleteItem)
	}
}

type createItemRequest struct {
	EncryptedData json.RawMessage `json:"encrypted_data" binding:"required"`
	Category      string          `json:"category"`
}

type updateItemRequest struct {
	EncryptedData json.RawMessage `json:"encrypted_data" binding:"required"`
	Category      *string         `json:"category"`
}

func (h *VaultHandler) createItem(c *gin.Context) {
	var req createItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apierror.ErrInvalidRequest)
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		handleError(c, apierror.ErrUnauthorized)
		return
	}

	newItem := &core.VaultItem{
		UserID:        userID.(uuid.UUID),
		EncryptedData: req.EncryptedData,
		Category:      req.Category,
	}

	createdItem, err := h.vaultService.CreateVaultItem(c.Request.Context(), newItem)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, createdItem)
}

func (h *VaultHandler) getItems(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		handleError(c, apierror.ErrUnauthorized)
		return
	}

	items, err := h.vaultService.GetVaultItems(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *VaultHandler) updateItem(c *gin.Context) {
	idParam := c.Param("id")
	itemID, err := uuid.Parse(idParam)
	if err != nil {
		handleError(c, apierror.New(http.StatusBadRequest, "Invalid item ID"))
		return
	}

	var req updateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apierror.ErrInvalidRequest)
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		handleError(c, apierror.ErrUnauthorized)
		return
	}

	itemToUpdate := &core.VaultItem{
		ID:            itemID,
		EncryptedData: req.EncryptedData,
	}
	if req.Category != nil {
		itemToUpdate.Category = *req.Category
	}

	updatedItem, err := h.vaultService.UpdateVaultItem(c.Request.Context(), itemToUpdate, userID.(uuid.UUID))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, updatedItem)
}

func (h *VaultHandler) deleteItem(c *gin.Context) {
	idParam := c.Param("id")
	itemID, err := uuid.Parse(idParam)
	if err != nil {
		handleError(c, apierror.New(http.StatusBadRequest, "Invalid item ID"))
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		handleError(c, apierror.ErrUnauthorized)
		return
	}

	err = h.vaultService.DeleteVaultItem(c.Request.Context(), itemID, userID.(uuid.UUID))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}