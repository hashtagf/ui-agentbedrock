# Requirements Document - Document Upload Feature

**Project**: UI AgentBedrock Test Interface
**Feature**: Document Upload
**Version**: 1.0
**Date**: 2025-12-17
**Branch**: main

---

## 1. Feature Overview

### 1.1 Purpose
เพิ่มความสามารถให้ผู้ใช้สามารถอัปโหลดเอกสาร (PDF, DOCX, TXT, etc.) เพื่อใช้เป็น context ในการสนทนากับ AgentBedrock

### 1.2 Goals
- ให้ผู้ใช้สามารถอัปโหลดไฟล์เอกสารได้ผ่าน chat interface
- แสดงไฟล์ที่อัปโหลดใน conversation history
- ส่งเนื้อหาจากไฟล์ไปให้ AgentBedrock เพื่อใช้เป็น context
- รองรับหลายประเภทไฟล์ (PDF, DOCX, TXT, images, etc.)

---

## 2. Functional Requirements

### FR-DU-001: File Upload UI
| ID | FR-DU-001 |
|----|-----------|
| **Title** | File Upload Interface |
| **Description** | ระบบต้องมี UI สำหรับอัปโหลดไฟล์ใน chat interface |
| **Priority** | High |
| **Acceptance Criteria** | <ul><li>มีปุ่มหรือ drag-and-drop area สำหรับอัปโหลดไฟล์</li><li>แสดง progress indicator ขณะอัปโหลด</li><li>แสดงชื่อไฟล์และขนาดไฟล์</li><li>รองรับการอัปโหลดหลายไฟล์พร้อมกัน</li><li>แสดง error message หากอัปโหลดไม่สำเร็จ</li></ul> |

### FR-DU-002: File Type Support
| ID | FR-DU-002 |
|----|-----------|
| **Title** | Supported File Types |
| **Description** | ระบบต้องรองรับไฟล์หลายประเภท |
| **Priority** | High |
| **Acceptance Criteria** | <ul><li>รองรับ PDF (.pdf)</li><li>รองรับ Microsoft Word (.docx, .doc)</li><li>รองรับ Text files (.txt, .md)</li><li>รองรับ Images (.png, .jpg, .jpeg) - optional</li><li>แสดง error หากไฟล์ไม่รองรับ</li><li>จำกัดขนาดไฟล์ (e.g., 10MB per file)</li></ul> |

### FR-DU-003: File Storage
| ID | FR-DU-003 |
|----|-----------|
| **Title** | File Storage Management |
| **Description** | ระบบต้องเก็บไฟล์ที่อัปโหลด |
| **Priority** | High |
| **Acceptance Criteria** | <ul><li>เก็บไฟล์ใน MongoDB GridFS หรือ local filesystem</li><li>เชื่อมโยงไฟล์กับ session</li><li>เก็บ metadata (ชื่อไฟล์, ขนาด, ประเภท, upload date)</li><li>สามารถลบไฟล์ได้เมื่อลบ session</li></ul> |

### FR-DU-004: Text Extraction
| ID | FR-DU-004 |
|----|-----------|
| **Title** | Document Text Extraction |
| **Description** | ระบบต้องสกัดข้อความจากเอกสารเพื่อส่งให้ AgentBedrock |
| **Priority** | High |
| **Acceptance Criteria** | <ul><li>สกัดข้อความจาก PDF ได้</li><li>สกัดข้อความจาก DOCX ได้</li><li>อ่านไฟล์ TXT ได้</li><li>แสดง error หากสกัดข้อความไม่สำเร็จ</li><li>แสดง preview ของข้อความที่สกัดได้ (optional)</li></ul> |

### FR-DU-005: Document Context Integration
| ID | FR-DU-005 |
|----|-----------|
| **Title** | Send Document Content to AgentBedrock |
| **Description** | ระบบต้องส่งเนื้อหาจากเอกสารไปให้ AgentBedrock เป็น context |
| **Priority** | High |
| **Acceptance Criteria** | <ul><li>รวมเนื้อหาจากเอกสารใน message context</li><li>ส่งพร้อมกับ user message ไปยัง AgentBedrock</li><li>แสดงใน conversation history ว่า message มีเอกสารแนบ</li><li>รองรับหลายเอกสารใน message เดียว</li></ul> |

### FR-DU-006: Document Display in Chat
| ID | FR-DU-006 |
|----|-----------|
| **Title** | Document Display in Message History |
| **Description** | แสดงเอกสารที่อัปโหลดใน conversation history |
| **Priority** | Medium |
| **Acceptance Criteria** | <ul><li>แสดงชื่อไฟล์และไอคอนประเภทไฟล์</li><li>แสดงขนาดไฟล์</li><li>สามารถดาวน์โหลดไฟล์เดิมได้ (optional)</li><li>แสดง preview สำหรับ images (optional)</li></ul> |

### FR-DU-007: File Size and Validation
| ID | FR-DU-007 |
|----|-----------|
| **Title** | File Validation and Size Limits |
| **Description** | ตรวจสอบและจำกัดขนาดไฟล์ |
| **Priority** | High |
| **Acceptance Criteria** | <ul><li>จำกัดขนาดไฟล์ต่อไฟล์ (e.g., 10MB)</li><li>จำกัดขนาดรวมต่อ message (e.g., 50MB)</li><li>ตรวจสอบประเภทไฟล์ก่อนอัปโหลด</li><li>แสดง error message ที่ชัดเจน</li></ul> |

---

## 3. Non-Functional Requirements

### NFR-DU-001: Performance
| ID | NFR-DU-001 |
|----|------------|
| **Title** | Upload Performance |
| **Description** | การอัปโหลดไฟล์ต้องเสร็จภายในเวลาที่เหมาะสม |
| **Priority** | Medium |
| **Acceptance Criteria** | <ul><li>ไฟล์ขนาด 1MB อัปโหลดเสร็จภายใน 5 วินาที</li><li>แสดง progress indicator ที่อัปเดตแบบ real-time</li></ul> |

### NFR-DU-002: Security
| ID | NFR-DU-002 |
|----|------------|
| **Title** | File Upload Security |
| **Description** | ต้องป้องกันการอัปโหลดไฟล์ที่เป็นอันตราย |
| **Priority** | High |
| **Acceptance Criteria** | <ul><li>ตรวจสอบประเภทไฟล์จาก content (ไม่ใช่แค่ extension)</li><li>สแกนไวรัส (optional - future enhancement)</li><li>จำกัดประเภทไฟล์ที่อนุญาต</li></ul> |

### NFR-DU-003: Storage Efficiency
| ID | NFR-DU-003 |
|----|------------|
| **Title** | Storage Management |
| **Description** | จัดการพื้นที่เก็บข้อมูลอย่างมีประสิทธิภาพ |
| **Priority** | Medium |
| **Acceptance Criteria** | <ul><li>ลบไฟล์อัตโนมัติเมื่อลบ session</li><li>มี cleanup job สำหรับไฟล์ที่ไม่ได้ใช้ (optional)</li></ul> |

---

## 4. Technical Requirements

### 4.1 File Storage Options

**Option A: MongoDB GridFS** (Recommended for small-medium files)
- Pros: Integrated with existing MongoDB, easy to manage
- Cons: Limited for very large files, may impact DB performance

**Option B: Local Filesystem**
- Pros: Simple, fast access
- Cons: Requires volume management in Docker, backup complexity

**Option C: AWS S3** (Future enhancement)
- Pros: Scalable, reliable, cost-effective for large files
- Cons: Additional AWS service, more complex setup

**Initial Choice**: MongoDB GridFS (can migrate to S3 later if needed)

### 4.2 Text Extraction Libraries

**For PDF:**
- Go: `github.com/gen2brain/go-fitz` or `github.com/pdfcpu/pdfcpu`
- Alternative: Use external service or AWS Textract

**For DOCX:**
- Go: `github.com/unidoc/unioffice` or `github.com/lukasjarosch/go-docx`

**For TXT:**
- Native Go file reading

**For Images (optional):**
- AWS Textract for OCR
- Or use AgentBedrock's native image support (if available)

### 4.3 API Endpoints (New)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/upload` | Upload file(s) |
| GET | `/api/files/:id` | Download file |
| DELETE | `/api/files/:id` | Delete file |
| GET | `/api/sessions/:id/files` | List files in session |

### 4.4 Data Model Changes

**New Model: Document**
```go
type Document struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    SessionID primitive.ObjectID `bson:"session_id" json:"sessionId"`
    MessageID primitive.ObjectID `bson:"message_id,omitempty" json:"messageId,omitempty"`
    Filename  string             `bson:"filename" json:"filename"`
    FileType  string             `bson:"file_type" json:"fileType"`
    FileSize  int64              `bson:"file_size" json:"fileSize"`
    Content   string             `bson:"content,omitempty" json:"content,omitempty"` // Extracted text
    GridFSID  primitive.ObjectID `bson:"gridfs_id,omitempty" json:"gridfsId,omitempty"`
    CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
}
```

**Updated Model: Message**
```go
type Message struct {
    // ... existing fields ...
    Documents []primitive.ObjectID `bson:"documents,omitempty" json:"documents,omitempty"` // NEW
}
```

**Updated Model: ChatRequest**
```go
type ChatRequest struct {
    SessionID string   `json:"sessionId" binding:"required"`
    Message   string   `json:"message" binding:"required"`
    DocumentIDs []string `json:"documentIds,omitempty"` // NEW
}
```

---

## 5. User Flow

### 5.1 Upload Document Flow
1. User clicks upload button or drags file into chat input area
2. File is validated (type, size)
3. File is uploaded to backend
4. Backend extracts text from document
5. Document metadata is saved to MongoDB
6. File is stored in GridFS
7. Frontend displays document in message input area
8. User sends message with document attached
9. Backend includes document content in message context
10. Message is sent to AgentBedrock with document context

### 5.2 Display Document in Chat Flow
1. Message with document is displayed in chat history
2. Document icon and filename are shown
3. User can click to view document details (optional)
4. Document content is included in message context when sending to agent

---

## 6. Constraints

- File size limit: 10MB per file, 50MB total per message
- Supported types: PDF, DOCX, DOC, TXT, MD (images optional)
- Storage: MongoDB GridFS (initial implementation)
- No file versioning (upload new file if needed)
- Files are session-scoped (deleted when session is deleted)

---

## 7. Out of Scope (Initial Version)

- Image OCR (can be added later)
- File editing or annotation
- File sharing between sessions
- Real-time collaborative document editing
- File versioning
- Advanced file preview (PDF viewer, etc.)
- AWS S3 storage (can be added later)

---

## 8. Open Questions

1. **Document Processing**: Should documents be processed immediately on upload, or only when message is sent?
   - **Recommendation**: Process on upload for better UX (show preview, validate content)

2. **Multiple Documents**: How should multiple documents be handled in one message?
   - **Recommendation**: Combine all document content into single context block

3. **Image Support**: Should images be supported initially?
   - **Recommendation**: Start with text documents (PDF, DOCX, TXT), add images later

4. **Document Persistence**: Should documents persist across messages in the same session?
   - **Recommendation**: Yes, documents are linked to session, can be referenced in multiple messages

5. **AgentBedrock Integration**: How should document content be sent to AgentBedrock?
   - **Recommendation**: Prepend document content to user message as context block

---

## 9. Implementation Priority

**Phase 1 (MVP)**:
- FR-DU-001: File Upload UI
- FR-DU-002: File Type Support (PDF, DOCX, TXT)
- FR-DU-003: File Storage (MongoDB GridFS)
- FR-DU-004: Text Extraction (basic)
- FR-DU-005: Document Context Integration
- FR-DU-007: File Validation

**Phase 2 (Enhancements)**:
- FR-DU-006: Document Display in Chat
- Image support
- Better preview
- Download functionality

---

## 10. Dependencies

- MongoDB GridFS driver for Go
- PDF extraction library
- DOCX extraction library
- Frontend file upload component
- File validation utilities

