package services

import (
	"bytes"
	"context"
	"fmt"
	"strings"
)

// ExtractionService handles text extraction from various document formats
type ExtractionService struct {
	// Future: Can add caching, rate limiting, etc.
}

func NewExtractionService() *ExtractionService {
	return &ExtractionService{}
}

// ExtractText extracts text content from a document based on file type
func (s *ExtractionService) ExtractText(ctx context.Context, fileType string, fileContent []byte) (string, error) {
	switch fileType {
	case "pdf":
		return s.extractFromPDF(ctx, fileContent)
	case "docx":
		return s.extractFromDOCX(ctx, fileContent)
	case "doc":
		return s.extractFromDOC(ctx, fileContent)
	case "txt", "md":
		return s.extractFromText(ctx, fileContent)
	default:
		return "", fmt.Errorf("unsupported file type: %s", fileType)
	}
}

// extractFromPDF extracts text from PDF files
// Note: This is a placeholder. In production, use a PDF library like go-fitz or pdfcpu
func (s *ExtractionService) extractFromPDF(ctx context.Context, content []byte) (string, error) {
	// TODO: Implement PDF extraction using github.com/gen2brain/go-fitz
	// For now, return a placeholder message
	// In production, this would use:
	// doc, err := fitz.NewFromMemory(content)
	// if err != nil { return "", err }
	// defer doc.Close()
	// var text strings.Builder
	// for i := 0; i < doc.NumPage(); i++ {
	//     pageText := doc.Text(i)
	//     text.WriteString(pageText)
	// }
	// return text.String(), nil

	return "", fmt.Errorf("PDF extraction not yet implemented. Please install PDF extraction library")
}

// extractFromDOCX extracts text from DOCX files
// Note: This is a placeholder. In production, use a DOCX library like unidoc/unioffice
func (s *ExtractionService) extractFromDOCX(ctx context.Context, content []byte) (string, error) {
	// TODO: Implement DOCX extraction using github.com/unidoc/unioffice
	// For now, return a placeholder message
	// In production, this would use:
	// docx, err := document.ReadDocxFromMemory(bytes.NewReader(content), int64(len(content)))
	// if err != nil { return "", err }
	// defer docx.Close()
	// return docx.GetContent(), nil

	return "", fmt.Errorf("DOCX extraction not yet implemented. Please install DOCX extraction library")
}

// extractFromDOC extracts text from DOC files (legacy Word format)
func (s *ExtractionService) extractFromDOC(ctx context.Context, content []byte) (string, error) {
	// DOC format is binary and complex. For MVP, we'll return an error suggesting conversion to DOCX
	return "", fmt.Errorf("DOC format is not supported. Please convert to DOCX format")
}

// extractFromText extracts text from plain text or markdown files
func (s *ExtractionService) extractFromText(ctx context.Context, content []byte) (string, error) {
	// Remove BOM if present
	content = bytes.TrimPrefix(content, []byte{0xEF, 0xBB, 0xBF})

	// Convert to string and clean up
	text := string(content)
	text = strings.TrimSpace(text)

	return text, nil
}
