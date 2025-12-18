package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Document represents an uploaded document file
type Document struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SessionID   primitive.ObjectID `bson:"session_id" json:"sessionId"`
	MessageID   primitive.ObjectID `bson:"message_id,omitempty" json:"messageId,omitempty"`
	Filename    string             `bson:"filename" json:"filename"`
	FileType    string             `bson:"file_type" json:"fileType"`                  // "pdf", "docx", "txt", "md", "xlsx", "xls"
	FileSize    int64              `bson:"file_size" json:"fileSize"`                  // bytes
	Content     string             `bson:"content,omitempty" json:"content,omitempty"` // Extracted text (not for Excel)
	GridFSID    primitive.ObjectID `bson:"gridfs_id,omitempty" json:"gridfsId,omitempty"`
	S3Key       string             `bson:"s3_key,omitempty" json:"s3Key,omitempty"`             // S3 object key for Excel files
	StorageType string             `bson:"storage_type,omitempty" json:"storageType,omitempty"` // "gridfs" or "s3"
	Confirmed   bool               `bson:"confirmed" json:"confirmed"`                          // True after S3 upload confirmed
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
}

// UploadResponse represents the response after successful file upload
type UploadResponse struct {
	DocumentID string `json:"documentId"`
	Filename   string `json:"filename"`
	FileType   string `json:"fileType"`
	FileSize   int64  `json:"fileSize"`
	Content    string `json:"content,omitempty"` // Extracted text preview (first 500 chars)
	S3Key      string `json:"s3Key,omitempty"`   // S3 key for Excel files
}
