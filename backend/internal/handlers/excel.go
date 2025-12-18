package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gin-gonic/gin"
	"github.com/ui-agentbedrock/backend/internal/models"
	"github.com/ui-agentbedrock/backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExcelHandler struct {
	lambdaClient   *lambda.Client
	lambdaFunction string
	documentRepo   *repository.DocumentRepository
}

func NewExcelHandler(cfg aws.Config, lambdaFunction string, documentRepo *repository.DocumentRepository) *ExcelHandler {
	return &ExcelHandler{
		lambdaClient:   lambda.NewFromConfig(cfg),
		lambdaFunction: lambdaFunction,
		documentRepo:   documentRepo,
	}
}

// PresignedURLRequest is the request body for getting a presigned URL
type PresignedURLRequest struct {
	SessionID string `json:"sessionId" binding:"required"`
	Filename  string `json:"filename" binding:"required"`
}

// PresignedURLResponse is the response with the presigned URL for S3 upload
type PresignedURLResponse struct {
	UploadURL  string `json:"uploadUrl"`
	FileKey    string `json:"fileKey"`
	BucketName string `json:"bucketName"`
	ExpiresIn  int    `json:"expiresIn"`
	DocumentID string `json:"documentId"` // Pre-created document ID for tracking
}

// LambdaRequest is the payload sent to the MCP Gateway Lambda
type LambdaRequest struct {
	Action     string                 `json:"action"`
	Parameters map[string]interface{} `json:"parameters"`
}

// BedrockLambdaResponse is the Bedrock Agent-formatted response from Lambda
type BedrockLambdaResponse struct {
	MessageVersion string `json:"messageVersion"`
	Response       struct {
		ActionGroup      string `json:"actionGroup"`
		Function         string `json:"function"`
		FunctionResponse struct {
			ResponseBody struct {
				TEXT struct {
					Body string `json:"body"`
				} `json:"TEXT"`
			} `json:"responseBody"`
		} `json:"functionResponse"`
	} `json:"response"`
}

// LambdaResponseData is the actual data inside the Bedrock response body
type LambdaResponseData struct {
	Action  string `json:"action"`
	Success bool   `json:"success"`
	Data    struct {
		UploadURL   string `json:"upload_url"`
		FileKey     string `json:"file_key"`
		Bucket      string `json:"bucket"`
		ExpiresIn   int    `json:"expires_in"`
		ContentType string `json:"content_type"`
		Filename    string `json:"filename"`
	} `json:"data"`
	Error string `json:"error,omitempty"`
}

// GetPresignedURL handles requests for S3 presigned upload URLs for Excel files
func (h *ExcelHandler) GetPresignedURL(c *gin.Context) {
	var req PresignedURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate session ID
	sessionID, err := primitive.ObjectIDFromHex(req.SessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sessionId"})
		return
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(req.Filename))
	if ext != ".xlsx" && ext != ".xls" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "only Excel files (.xlsx, .xls) can use presigned upload",
		})
		return
	}

	// Determine content type
	contentType := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	if ext == ".xls" {
		contentType = "application/vnd.ms-excel"
	}

	// Call Lambda to generate presigned URL
	lambdaPayload := LambdaRequest{
		Action: "generate_presigned_upload_url",
		Parameters: map[string]interface{}{
			"filename":     req.Filename,
			"content_type": contentType,
		},
	}

	payloadBytes, err := json.Marshal(lambdaPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to prepare request"})
		return
	}

	// Invoke Lambda
	result, err := h.lambdaClient.Invoke(c.Request.Context(), &lambda.InvokeInput{
		FunctionName: aws.String(h.lambdaFunction),
		Payload:      payloadBytes,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to call Lambda: %v", err)})
		return
	}

	// Parse Bedrock-formatted Lambda response
	var bedrockResp BedrockLambdaResponse
	if err := json.Unmarshal(result.Payload, &bedrockResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to parse Lambda response: %v", err)})
		return
	}

	// Extract the actual response data from the Bedrock response body
	var lambdaData LambdaResponseData
	if err := json.Unmarshal([]byte(bedrockResp.Response.FunctionResponse.ResponseBody.TEXT.Body), &lambdaData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to parse response body: %v", err)})
		return
	}

	if !lambdaData.Success {
		c.JSON(http.StatusInternalServerError, gin.H{"error": lambdaData.Error})
		return
	}

	// Pre-create document record in MongoDB for tracking
	doc := &models.Document{
		SessionID:   sessionID,
		Filename:    req.Filename,
		FileType:    strings.TrimPrefix(ext, "."),
		S3Key:       lambdaData.Data.FileKey,
		StorageType: "s3",
	}

	if err := h.documentRepo.CreateDocument(c.Request.Context(), doc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create document record"})
		return
	}

	// Return presigned URL info
	c.JSON(http.StatusOK, PresignedURLResponse{
		UploadURL:  lambdaData.Data.UploadURL,
		FileKey:    lambdaData.Data.FileKey,
		BucketName: lambdaData.Data.Bucket,
		ExpiresIn:  lambdaData.Data.ExpiresIn,
		DocumentID: doc.ID.Hex(),
	})
}

// ConfirmExcelUpload is called after successful S3 upload to update document status
func (h *ExcelHandler) ConfirmExcelUpload(c *gin.Context) {
	documentIDStr := c.Param("id")
	documentID, err := primitive.ObjectIDFromHex(documentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document id"})
		return
	}

	var req struct {
		FileSize int64 `json:"fileSize"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update document with file size and confirm upload
	if err := h.documentRepo.ConfirmS3Upload(c.Request.Context(), documentID, req.FileSize); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to confirm upload"})
		return
	}

	// Get updated document
	doc, err := h.documentRepo.GetDocument(c.Request.Context(), documentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get document"})
		return
	}

	c.JSON(http.StatusOK, models.UploadResponse{
		DocumentID: doc.ID.Hex(),
		Filename:   doc.Filename,
		FileType:   doc.FileType,
		FileSize:   doc.FileSize,
		S3Key:      doc.S3Key,
	})
}

// Helper function to check if Lambda is configured
func (h *ExcelHandler) IsConfigured() bool {
	return h.lambdaFunction != ""
}

// GetExcelFile returns info about an Excel file (for download, we redirect to S3)
func (h *ExcelHandler) GetExcelFileInfo(ctx context.Context, documentID primitive.ObjectID) (*models.Document, error) {
	return h.documentRepo.GetDocument(ctx, documentID)
}
