package repository

import (
	"context"
	"io"
	"time"

	"github.com/ui-agentbedrock/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DocumentRepository struct {
	documents *mongo.Collection
	bucket    *gridfs.Bucket
}

func NewDocumentRepository(db *mongo.Database) *DocumentRepository {
	bucket, err := gridfs.NewBucket(db, options.GridFSBucket().SetName("documents"))
	if err != nil {
		panic("Failed to create GridFS bucket: " + err.Error())
	}

	return &DocumentRepository{
		documents: db.Collection("documents"),
		bucket:    bucket,
	}
}

// SaveDocument saves document metadata and file to GridFS
func (r *DocumentRepository) SaveDocument(ctx context.Context, doc *models.Document, fileReader io.Reader) error {
	doc.ID = primitive.NewObjectID()
	doc.CreatedAt = time.Now()

	// Upload file to GridFS
	uploadStream, err := r.bucket.OpenUploadStreamWithID(doc.ID, doc.Filename)
	if err != nil {
		return err
	}
	defer uploadStream.Close()

	_, err = io.Copy(uploadStream, fileReader)
	if err != nil {
		return err
	}

	// Save document metadata
	doc.GridFSID = doc.ID
	_, err = r.documents.InsertOne(ctx, doc)
	return err
}

// GetDocument retrieves document metadata by ID
func (r *DocumentRepository) GetDocument(ctx context.Context, id primitive.ObjectID) (*models.Document, error) {
	var doc models.Document
	err := r.documents.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// GetDocumentsBySession retrieves all documents for a session
func (r *DocumentRepository) GetDocumentsBySession(ctx context.Context, sessionID primitive.ObjectID) ([]models.Document, error) {
	cursor, err := r.documents.Find(ctx, bson.M{"session_id": sessionID}, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var documents []models.Document
	if err := cursor.All(ctx, &documents); err != nil {
		return nil, err
	}

	if documents == nil {
		documents = []models.Document{}
	}
	return documents, nil
}

// GetDocumentsByIDs retrieves multiple documents by their IDs
func (r *DocumentRepository) GetDocumentsByIDs(ctx context.Context, ids []primitive.ObjectID) ([]models.Document, error) {
	if len(ids) == 0 {
		return []models.Document{}, nil
	}

	cursor, err := r.documents.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var documents []models.Document
	if err := cursor.All(ctx, &documents); err != nil {
		return nil, err
	}

	if documents == nil {
		documents = []models.Document{}
	}
	return documents, nil
}

// DownloadFile retrieves file content from GridFS
func (r *DocumentRepository) DownloadFile(ctx context.Context, id primitive.ObjectID) (io.ReadCloser, error) {
	downloadStream, err := r.bucket.OpenDownloadStream(id)
	if err != nil {
		return nil, err
	}
	return downloadStream, nil
}

// DeleteDocument deletes document metadata and file from GridFS
func (r *DocumentRepository) DeleteDocument(ctx context.Context, id primitive.ObjectID) error {
	// Delete from GridFS
	if err := r.bucket.Delete(id); err != nil {
		// Ignore error if file doesn't exist in GridFS
		if err != gridfs.ErrFileNotFound {
			return err
		}
	}

	// Delete metadata
	_, err := r.documents.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// DeleteDocumentsBySession deletes all documents for a session
func (r *DocumentRepository) DeleteDocumentsBySession(ctx context.Context, sessionID primitive.ObjectID) error {
	// Get all documents for the session
	documents, err := r.GetDocumentsBySession(ctx, sessionID)
	if err != nil {
		return err
	}

	// Delete each document
	for _, doc := range documents {
		if err := r.DeleteDocument(ctx, doc.ID); err != nil {
			return err
		}
	}

	return nil
}

// UpdateDocumentContent updates the extracted text content of a document
func (r *DocumentRepository) UpdateDocumentContent(ctx context.Context, id primitive.ObjectID, content string) error {
	_, err := r.documents.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"content": content}},
	)
	return err
}

// CreateDocument creates a document record without file content (for S3 uploads)
func (r *DocumentRepository) CreateDocument(ctx context.Context, doc *models.Document) error {
	doc.ID = primitive.NewObjectID()
	doc.CreatedAt = time.Now()
	doc.Confirmed = false

	_, err := r.documents.InsertOne(ctx, doc)
	return err
}

// ConfirmS3Upload updates a document after successful S3 upload
func (r *DocumentRepository) ConfirmS3Upload(ctx context.Context, id primitive.ObjectID, fileSize int64) error {
	_, err := r.documents.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{
			"file_size": fileSize,
			"confirmed": true,
		}},
	)
	return err
}
