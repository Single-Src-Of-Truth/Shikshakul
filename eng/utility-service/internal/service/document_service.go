package service

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Single-Src-Of-Truth/Shikshakul/eng/utility-service/internal/domain"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type DocumentService struct {
	storage domain.StorageProvider
	logger  *zap.Logger
}

func NewDocumentService(provider domain.StorageProvider, logger *zap.Logger) *DocumentService {
	return &DocumentService{
		storage: provider,
		logger:  logger,
	}
}

func (s *DocumentService) RequestUpload(ctx context.Context, tenantID, category, filename, contentType string) (string, string, error) {
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = ".bin"
	}

	newFileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	objectKey := fmt.Sprintf("%s/temp/%s/%s", tenantID, category, newFileName)

	url, err := s.storage.GenerateUploadURL(ctx, objectKey, contentType)
	if err != nil {
		s.logger.Error("Failed to generate upload URL", zap.Error(err))
		return "", "", fmt.Errorf("%s: %v", domain.ErrProviderFailed, err)
	}

	return url, objectKey, nil
}

func (s *DocumentService) ApproveDocument(ctx context.Context, sourceKey string) (string, error) {
	if !strings.Contains(sourceKey, "/temp/") {
		return "", fmt.Errorf("invalid source key: document is not in temporary storage")
	}

	destKey := strings.Replace(sourceKey, "/temp/", "/permanent/", 1)

	err := s.storage.MoveObject(ctx, sourceKey, destKey)
	if err != nil {
		s.logger.Error("Failed to move document", zap.String("source", sourceKey), zap.Error(err))
		return "", fmt.Errorf("%s: %v", domain.ErrProviderFailed, err)
	}

	return destKey, nil
}

func (s *DocumentService) GetAccessURL(ctx context.Context, objectKey string, viewOnly bool) (string, error) {
	url, err := s.storage.GenerateDownloadURL(ctx, objectKey, viewOnly)
	if err != nil {
		s.logger.Error("Failed to generate access URL", zap.String("key", objectKey), zap.Error(err))
		return "", fmt.Errorf("%s: %v", domain.ErrProviderFailed, err)
	}
	return url, nil
}

func (s *DocumentService) SoftDeleteDocument(ctx context.Context, sourceKey string) (string, error) {
	var destKey string

	if strings.Contains(sourceKey, "/permanent/") {
		destKey = strings.Replace(sourceKey, "/permanent/", "/trash/", 1)
	} else if strings.Contains(sourceKey, "/temp/") {
		destKey = strings.Replace(sourceKey, "/temp/", "/trash/", 1)
	} else {
		return "", fmt.Errorf("invalid source key: document must be in temp or permanent folders")
	}

	err := s.storage.MoveObject(ctx, sourceKey, destKey)
	if err != nil {
		s.logger.Error("Failed to soft-delete document", zap.String("source", sourceKey), zap.Error(err))
		return "", fmt.Errorf("%s: %v", domain.ErrProviderFailed, err)
	}

	return destKey, nil
}
