package handlers

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ui-agentbedrock/backend/internal/models"
	"github.com/ui-agentbedrock/backend/internal/repository"
	"github.com/ui-agentbedrock/backend/internal/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	MaxFileSize      = 10 * 1024 * 1024 // 10MB per file
	MaxTotalFileSize = 50 * 1024 * 1024 // 50MB total per message
)

var allowedMimeTypes = map[string]string{
	"application/pdf": "pdf",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": "docx",
	"application/msword": "doc",
	"text/plain":         "txt",
	"text/markdown":      "md",
	// Excel formats - handled differently (presigned S3 upload)
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": "xlsx",
	"application/vnd.ms-excel": "xls",
}

// ExcelMimeTypes for files that should be uploaded to S3 instead of GridFS
var ExcelMimeTypes = map[string]bool{
	"xlsx": true,
	"xls":  true,
}

type UploadHandler struct {
	documentRepo   *repository.DocumentRepository
	extractService *services.ExtractionService
}

func NewUploadHandler(documentRepo *repository.DocumentRepository, extractService *services.ExtractionService) *UploadHandler {
	return &UploadHandler{
		documentRepo:   documentRepo,
		extractService: extractService,
	}
}

// UploadFile handles file upload
func (h *UploadHandler) UploadFile(c *gin.Context) {
	// Get session ID from form
	sessionIDStr := c.PostForm("sessionId")
	if sessionIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sessionId is required"})
		return
	}

	sessionID, err := primitive.ObjectIDFromHex(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sessionId"})
		return
	}

	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	// Validate file size
	if file.Size > MaxFileSize {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{
			"error": fmt.Sprintf("file size exceeds maximum allowed size of %d bytes", MaxFileSize),
		})
		return
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer src.Close()

	// Read file content for validation
	fileContent := make([]byte, file.Size)
	if _, err := src.Read(fileContent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}

	// Detect MIME type
	mimeType := http.DetectContentType(fileContent)
	fileType, ok := allowedMimeTypes[mimeType]
	if !ok {
		// Try to get from extension as fallback
		ext := strings.ToLower(filepath.Ext(file.Filename))
		switch ext {
		case ".pdf":
			fileType = "pdf"
		case ".docx":
			fileType = "docx"
		case ".doc":
			fileType = "doc"
		case ".txt":
			fileType = "txt"
		case ".md":
			fileType = "md"
		case ".xlsx":
			fileType = "xlsx"
		case ".xls":
			fileType = "xls"
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("unsupported file type: %s. Allowed types: PDF, DOCX, DOC, TXT, MD, XLSX, XLS", mimeType),
			})
			return
		}
	}

	// Create document model
	doc := &models.Document{
		SessionID: sessionID,
		Filename:  file.Filename,
		FileType:  fileType,
		FileSize:  file.Size,
	}

	// Save document to GridFS
	fileReader := strings.NewReader(string(fileContent))
	if err := h.documentRepo.SaveDocument(c.Request.Context(), doc, fileReader); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save document"})
		return
	}

	// Extract text content
	content, extractErr := h.extractService.ExtractText(c.Request.Context(), fileType, fileContent)
	if extractErr != nil {
		// Log error but don't fail upload
		fmt.Printf("Warning: Failed to extract text from document %s: %v\n", doc.ID.Hex(), extractErr)
		// For TXT/MD files, extraction should always work, so this is unexpected
		if fileType == "txt" || fileType == "md" {
			fmt.Printf("Error: Text extraction failed for text file, this should not happen\n")
		}
	} else {
		// Update document with extracted content
		doc.Content = content
		if err := h.documentRepo.UpdateDocumentContent(c.Request.Context(), doc.ID, content); err != nil {
			fmt.Printf("Warning: Failed to update document content: %v\n", err)
		}
	}

	// Prepare response
	preview := content
	if len(preview) > 500 {
		preview = preview[:500] + "..."
	}

	response := models.UploadResponse{
		DocumentID: doc.ID.Hex(),
		Filename:   doc.Filename,
		FileType:   doc.FileType,
		FileSize:   doc.FileSize,
		Content:    preview,
	}

	c.JSON(http.StatusOK, response)
}

// DownloadFile handles file download
func (h *UploadHandler) DownloadFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document id"})
		return
	}

	// Get document metadata
	doc, err := h.documentRepo.GetDocument(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}

	// Download file from GridFS
	fileStream, err := h.documentRepo.DownloadFile(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to download file"})
		return
	}
	defer fileStream.Close()

	// Set headers
	mimeType := getMimeType(doc.FileType)
	c.Header("Content-Type", mimeType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", doc.Filename))
	c.Header("Content-Length", fmt.Sprintf("%d", doc.FileSize))

	// Stream file
	c.Stream(func(w io.Writer) bool {
		_, err := io.Copy(w, fileStream)
		return err == nil
	})
}

// DeleteFile handles file deletion
func (h *UploadHandler) DeleteFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document id"})
		return
	}

	if err := h.documentRepo.DeleteDocument(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete document"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// GetSessionDocuments retrieves all documents for a session
func (h *UploadHandler) GetSessionDocuments(c *gin.Context) {
	sessionIDStr := c.Param("id")
	sessionID, err := primitive.ObjectIDFromHex(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}

	documents, err := h.documentRepo.GetDocumentsBySession(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get documents"})
		return
	}

	c.JSON(http.StatusOK, documents)
}

// Helper function to get MIME type from file type
func getMimeType(fileType string) string {
	switch fileType {
	case "pdf":
		return "application/pdf"
	case "docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case "doc":
		return "application/msword"
	case "txt":
		return "text/plain"
	case "md":
		return "text/markdown"
	case "xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case "xls":
		return "application/vnd.ms-excel"
	default:
		return "application/octet-stream"
	}
}

// IsExcelFile checks if the file type is an Excel file
func IsExcelFile(fileType string) bool {
	return ExcelMimeTypes[fileType]
}
