interface Document {
  id: string
  filename: string
  fileType: string
  fileSize: number
  content?: string
  s3Key?: string
  storageType?: 'gridfs' | 's3'
  createdAt: string
}

interface UploadedDocument {
  documentId: string
  filename: string
  fileType: string
  fileSize: number
  content?: string
  s3Key?: string
}

interface PresignedURLResponse {
  uploadUrl: string
  fileKey: string
  bucketName: string
  expiresIn: number
  documentId: string
}

// Excel file extensions that should use presigned S3 upload
const EXCEL_EXTENSIONS = ['.xlsx', '.xls']

export function useDocumentUpload() {
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase
  
  const uploadedDocuments = useState<UploadedDocument[]>('uploadedDocuments', () => [])
  const isUploading = useState<boolean>('isUploading', () => false)
  const uploadError = useState<string | null>('uploadError', () => null)
  const uploadProgress = useState<number>('uploadProgress', () => 0)

  const uploadFile = async (file: File, sessionId: string): Promise<UploadedDocument | null> => {
    if (!sessionId) {
      uploadError.value = 'No session selected'
      return null
    }

    // Validate file size (10MB)
    const maxSize = 10 * 1024 * 1024
    if (file.size > maxSize) {
      uploadError.value = `File size exceeds maximum allowed size of ${(maxSize / 1024 / 1024).toFixed(0)}MB`
      return null
    }

    // Validate file type
    const allowedTypes = [
      'application/pdf', 
      'application/vnd.openxmlformats-officedocument.wordprocessingml.document', 
      'application/msword', 
      'text/plain', 
      'text/markdown',
      'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
      'application/vnd.ms-excel'
    ]
    const allowedExtensions = ['.pdf', '.docx', '.doc', '.txt', '.md', '.xlsx', '.xls']
    const fileExtension = '.' + file.name.split('.').pop()?.toLowerCase()
    
    if (!allowedTypes.includes(file.type) && !allowedExtensions.includes(fileExtension)) {
      uploadError.value = 'Unsupported file type. Allowed types: PDF, DOCX, DOC, TXT, MD, XLSX, XLS'
      return null
    }

    // Check if this is an Excel file - use presigned S3 upload
    if (EXCEL_EXTENSIONS.includes(fileExtension)) {
      return await uploadExcelFile(file, sessionId)
    }

    isUploading.value = true
    uploadError.value = null
    uploadProgress.value = 0

    try {
      const formData = new FormData()
      formData.append('sessionId', sessionId)
      formData.append('file', file)

      const xhr = new XMLHttpRequest()

      // Track upload progress
      xhr.upload.addEventListener('progress', (e) => {
        if (e.lengthComputable) {
          uploadProgress.value = Math.round((e.loaded / e.total) * 100)
        }
      })

      const uploadPromise = new Promise<UploadedDocument>((resolve, reject) => {
        xhr.addEventListener('load', () => {
          if (xhr.status === 200) {
            const response = JSON.parse(xhr.responseText)
            resolve(response)
          } else {
            const error = JSON.parse(xhr.responseText)
            reject(new Error(error.error || 'Upload failed'))
          }
        })

        xhr.addEventListener('error', () => {
          reject(new Error('Network error during upload'))
        })

        xhr.addEventListener('abort', () => {
          reject(new Error('Upload cancelled'))
        })
      })

      xhr.open('POST', `${apiBase}/api/upload`)
      xhr.send(formData)

      const result = await uploadPromise
      uploadedDocuments.value = [...uploadedDocuments.value, result]
      uploadProgress.value = 100
      
      return result
    } catch (error: any) {
      uploadError.value = error.message || 'Failed to upload file'
      return null
    } finally {
      isUploading.value = false
      setTimeout(() => {
        uploadProgress.value = 0
      }, 1000)
    }
  }

  // Upload Excel file using presigned S3 URL
  const uploadExcelFile = async (file: File, sessionId: string): Promise<UploadedDocument | null> => {
    isUploading.value = true
    uploadError.value = null
    uploadProgress.value = 0

    try {
      // Step 1: Get presigned URL from backend
      const presignResponse = await fetch(`${apiBase}/api/excel/presign`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          sessionId,
          filename: file.name
        })
      })

      if (!presignResponse.ok) {
        const error = await presignResponse.json()
        throw new Error(error.error || 'Failed to get upload URL')
      }

      const presignData: PresignedURLResponse = await presignResponse.json()
      uploadProgress.value = 20

      // Step 2: Upload file directly to S3 using presigned PUT URL
      const s3Response = await fetch(presignData.uploadUrl, {
        method: 'PUT',
        body: file,
        headers: {
          'Content-Type': file.type || 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
        }
      })

      if (!s3Response.ok) {
        throw new Error('Failed to upload file to S3')
      }
      uploadProgress.value = 80

      // Step 3: Confirm upload with backend
      const confirmResponse = await fetch(`${apiBase}/api/excel/confirm/${presignData.documentId}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          fileSize: file.size
        })
      })

      if (!confirmResponse.ok) {
        const error = await confirmResponse.json()
        throw new Error(error.error || 'Failed to confirm upload')
      }

      const result: UploadedDocument = await confirmResponse.json()
      uploadedDocuments.value = [...uploadedDocuments.value, result]
      uploadProgress.value = 100
      
      return result
    } catch (error: any) {
      uploadError.value = error.message || 'Failed to upload Excel file'
      return null
    } finally {
      isUploading.value = false
      setTimeout(() => {
        uploadProgress.value = 0
      }, 1000)
    }
  }

  const removeDocument = (documentId: string) => {
    uploadedDocuments.value = uploadedDocuments.value.filter(doc => doc.documentId !== documentId)
  }

  const clearDocuments = () => {
    uploadedDocuments.value = []
    uploadError.value = null
  }

  const getDocumentIds = (): string[] => {
    return uploadedDocuments.value.map(doc => doc.documentId)
  }

  return {
    uploadedDocuments,
    isUploading,
    uploadError,
    uploadProgress,
    uploadFile,
    removeDocument,
    clearDocuments,
    getDocumentIds,
  }
}

