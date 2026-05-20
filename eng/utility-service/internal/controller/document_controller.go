package controller

import (
	"net/http"
	"time"

	"github.com/Single-Src-Of-Truth/Shikshakul/eng/utility-service/internal/service"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/utility-service/pkg/dto"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/utility-service/pkg/response"
	"github.com/gin-gonic/gin"
)

type DocumentController struct {
	docService *service.DocumentService
}

func NewDocumentController(docService *service.DocumentService) *DocumentController {
	return &DocumentController{
		docService: docService,
	}
}

func (c *DocumentController) SignUpload(ctx *gin.Context) {
	var req dto.SignUploadRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Error("Invalid request payload", err.Error()))
		return
	}

	url, objectKey, err := c.docService.RequestUpload(ctx.Request.Context(), req.TenantID, req.Category, req.FileName, req.ContentType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Error("Failed to sign upload request", err.Error()))
		return
	}

	resp := dto.SignUploadResponse{
		UploadURL: url,
		ObjectKey: objectKey,
		ExpiresIn: 900,
	}

	ctx.JSON(http.StatusOK, response.Success("Upload URL generated successfully", resp))
}

func (c *DocumentController) MoveDocument(ctx *gin.Context) {
	var req dto.MoveDocumentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Error("Invalid request payload", err.Error()))
		return
	}

	newKey, err := c.docService.ApproveDocument(ctx.Request.Context(), req.SourceKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Error("Failed to move document", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.Success("Document moved successfully", gin.H{"new_key": newKey}))
}

func (c *DocumentController) GetAccessURL(ctx *gin.Context) {
	objectKey := ctx.Query("key")
	if objectKey == "" {
		ctx.JSON(http.StatusBadRequest, response.Error("Missing document key", "The 'key' query parameter is required"))
		return
	}

	isView := ctx.Query("view") == "true"

	url, err := c.docService.GetAccessURL(ctx.Request.Context(), objectKey, isView)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Error("Failed to generate access URL", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.Success("URL generated successfully", gin.H{"url": url}))
}

func (c *DocumentController) SoftDelete(ctx *gin.Context) {
	var req dto.SoftDeleteDocumentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Error("Invalid request payload", err.Error()))
		return
	}

	newKey, err := c.docService.SoftDeleteDocument(ctx.Request.Context(), req.SourceKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Error("Failed to soft delete document", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.Success("Document moved to trash successfully", gin.H{"new_key": newKey}))
}

func (c *DocumentController) Ping(ctx *gin.Context) {
	healthStatuses := c.docService.HealthCheck(ctx.Request.Context())

	isFullyHealthy := true
	for _, state := range healthStatuses {
		if state != "UP" {
			isFullyHealthy = false
			break
		}
	}

	responsePayload := gin.H{
		"success":    isFullyHealthy,
		"timestamp":  time.Now().Format(time.RFC3339),
		"components": healthStatuses,
	}

	if isFullyHealthy {
		responsePayload["status"] = "UP"
		ctx.JSON(http.StatusOK, responsePayload)
	} else {
		responsePayload["status"] = "DOWN"
		ctx.JSON(http.StatusServiceUnavailable, responsePayload)
	}
}
