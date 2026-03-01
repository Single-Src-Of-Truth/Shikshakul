package dto

type SignUploadRequest struct {
	TenantID    string `json:"tenant_id" binding:"required"`
	Category    string `json:"category" binding:"required"`
	FileName    string `json:"file_name" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
}

type SignUploadResponse struct {
	UploadURL string `json:"upload_url"`
	ObjectKey string `json:"object_key"`
	ExpiresIn int    `json:"expires_in"`
}

type MoveDocumentRequest struct {
	SourceKey string `json:"source_key" binding:"required"`
}

type SoftDeleteDocumentRequest struct {
	SourceKey string `json:"source_key" binding:"required"`
}
