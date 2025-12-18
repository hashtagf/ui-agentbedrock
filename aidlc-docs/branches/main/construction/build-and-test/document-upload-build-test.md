# Build and Test Instructions - Document Upload Feature

**Project**: UI AgentBedrock Test Interface
**Feature**: Document Upload
**Date**: 2025-12-17
**Branch**: main

---

## Build Instructions

### Prerequisites

1. **MongoDB**: Ensure MongoDB is running (via Docker Compose or local installation)
2. **Go**: Version 1.23 or higher
3. **Node.js**: Version 18 or higher
4. **Docker & Docker Compose**: For containerized deployment

### Backend Build

```bash
cd backend
go mod tidy
go build -o server ./cmd/server
```

### Frontend Build

```bash
cd frontend
npm install
npm run build
```

### Docker Build

```bash
# Build all services
docker-compose build

# Or build individually
docker-compose build backend
docker-compose build frontend
```

---

## Test Instructions

### Unit Tests

#### Backend Tests

**Note**: Create test files as needed:

```bash
cd backend

# Test Document Repository
go test ./internal/repository -run TestDocumentRepository

# Test Upload Handler
go test ./internal/handlers -run TestUploadHandler

# Test Extraction Service
go test ./internal/services -run TestExtractionService
```

#### Frontend Tests

```bash
cd frontend

# Run unit tests (if configured)
npm test

# Or run with Vitest
npm run test:unit
```

### Integration Tests

#### 1. Test File Upload Endpoint

```bash
# Start backend server
cd backend
go run ./cmd/server

# In another terminal, test upload
curl -X POST http://localhost:8081/api/upload \
  -F "sessionId=YOUR_SESSION_ID" \
  -F "file=@test.txt"
```

**Expected Response**:
```json
{
  "documentId": "...",
  "filename": "test.txt",
  "fileType": "txt",
  "fileSize": 1234,
  "content": "File content preview..."
}
```

#### 2. Test Document Download

```bash
curl http://localhost:8081/api/files/DOCUMENT_ID \
  --output downloaded_file.txt
```

#### 3. Test Document List

```bash
curl http://localhost:8081/api/sessions/SESSION_ID/documents
```

#### 4. Test Chat with Documents

```bash
curl -X POST http://localhost:8081/api/chat/stream \
  -H "Content-Type: application/json" \
  -d '{
    "sessionId": "SESSION_ID",
    "message": "What is in the document?",
    "documentIds": ["DOCUMENT_ID_1", "DOCUMENT_ID_2"]
  }'
```

### End-to-End Tests

#### Manual E2E Test Flow

1. **Start Services**:
   ```bash
   docker-compose up
   ```

2. **Access Frontend**: Open http://localhost:3000

3. **Create Session**: Click "New Chat" button

4. **Upload Document**:
   - Click paperclip icon in input area
   - Select a TXT or MD file (PDF/DOCX require libraries)
   - Verify upload progress indicator
   - Verify document appears in document list

5. **Send Message with Document**:
   - Type a message
   - Click send
   - Verify message includes document context
   - Verify AI response references document content

6. **Verify Document in History**:
   - Check message history
   - Verify document indicator shows on message
   - Verify document count is correct

7. **Test Multiple Documents**:
   - Upload multiple files
   - Send message with all documents
   - Verify all documents are included

8. **Test Error Handling**:
   - Try uploading file > 10MB (should fail)
   - Try uploading unsupported file type (should fail)
   - Verify error messages are displayed

### Test Cases

#### File Upload Tests

| Test Case | Expected Result |
|-----------|----------------|
| Upload TXT file | ✅ Success, text extracted |
| Upload MD file | ✅ Success, text extracted |
| Upload PDF file | ⚠️ Upload succeeds, extraction placeholder (needs library) |
| Upload DOCX file | ⚠️ Upload succeeds, extraction placeholder (needs library) |
| Upload file > 10MB | ❌ Error: File too large |
| Upload unsupported type | ❌ Error: Unsupported file type |
| Upload multiple files | ✅ All files uploaded successfully |
| Drag and drop file | ✅ File uploaded via drag-drop |

#### Document Integration Tests

| Test Case | Expected Result |
|-----------|----------------|
| Send message with document | ✅ Document content included in context |
| Send message with multiple documents | ✅ All documents combined in context |
| Send message without document | ✅ Normal message flow |
| View message with documents | ✅ Document indicator displayed |
| Remove document before sending | ✅ Document removed from message |

#### Storage Tests

| Test Case | Expected Result |
|-----------|----------------|
| Document stored in GridFS | ✅ File accessible via download endpoint |
| Document metadata saved | ✅ Metadata queryable |
| Delete document | ✅ File and metadata deleted |
| Delete session | ✅ All session documents deleted |

---

## Known Limitations

1. **PDF Extraction**: Not yet implemented (requires `github.com/gen2brain/go-fitz` or similar)
2. **DOCX Extraction**: Not yet implemented (requires `github.com/unidoc/unioffice` or similar)
3. **DOC Format**: Not supported (legacy format, suggest conversion to DOCX)

---

## Installation of Extraction Libraries (Optional)

### For PDF Extraction

```bash
cd backend
go get github.com/gen2brain/go-fitz
```

Then update `services/extraction.go` to use the library.

### For DOCX Extraction

```bash
cd backend
go get github.com/unidoc/unioffice/document
```

Then update `services/extraction.go` to use the library.

---

## Troubleshooting

### Backend Issues

**Error: "Failed to create GridFS bucket"**
- Ensure MongoDB is running
- Check MongoDB connection string
- Verify database permissions

**Error: "Failed to extract text"**
- For PDF/DOCX: Libraries not installed (expected for MVP)
- For TXT/MD: Check file encoding (should be UTF-8)

**Error: "File too large"**
- Check file size (max 10MB per file)
- Check total size limit (50MB per message)

### Frontend Issues

**Error: "Cannot find name 'useDocumentUpload'"**
- Restart TypeScript server
- Verify `useDocumentUpload.ts` exists in `composables/` directory
- Check Nuxt auto-imports configuration

**Upload not working**
- Check browser console for errors
- Verify API endpoint is accessible
- Check CORS configuration

**Documents not displaying**
- Check browser console for errors
- Verify document IDs are being sent correctly
- Check MessageList component updates

---

## Performance Testing

### Upload Performance

- **Target**: 1MB file uploads in < 5 seconds
- **Test**: Upload various file sizes and measure time
- **Monitor**: Network latency, server processing time

### Extraction Performance

- **Target**: TXT/MD extraction in < 1 second
- **Test**: Extract text from files of various sizes
- **Monitor**: Processing time, memory usage

### Storage Performance

- **Target**: GridFS operations in < 2 seconds
- **Test**: Upload, download, delete operations
- **Monitor**: MongoDB performance, disk I/O

---

## Security Testing

1. **File Type Validation**: Test with malicious file types
2. **File Size Limits**: Test with oversized files
3. **Content Validation**: Test with corrupted files
4. **Session Scoping**: Verify documents only accessible within session
5. **Access Control**: Test unauthorized document access

---

## Next Steps

1. Install PDF extraction library (go-fitz)
2. Install DOCX extraction library (unidoc/unioffice)
3. Implement full extraction for PDF and DOCX
4. Add image support (optional)
5. Add document preview (optional)
6. Migrate to S3 storage (optional, for scalability)

---

## Test Checklist

- [ ] Backend builds successfully
- [ ] Frontend builds successfully
- [ ] Docker Compose starts all services
- [ ] File upload endpoint works
- [ ] TXT file extraction works
- [ ] MD file extraction works
- [ ] Document storage in GridFS works
- [ ] Document download works
- [ ] Document deletion works
- [ ] Chat with documents works
- [ ] Document display in UI works
- [ ] Error handling works
- [ ] File validation works
- [ ] Multiple file upload works
- [ ] Drag and drop works
- [ ] Progress indicator works

---

## Notes

- PDF and DOCX extraction are placeholders and will return errors until libraries are installed
- TXT and MD extraction is fully functional
- All file storage uses MongoDB GridFS
- Documents are session-scoped and deleted with sessions
- Maximum file size: 10MB per file, 50MB total per message

