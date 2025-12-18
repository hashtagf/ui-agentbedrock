interface Document {
  id: string
  filename: string
  fileType: string
  fileSize: number
  content?: string
  createdAt: string
}

interface UploadedDocument {
  documentId: string
  filename: string
  fileType: string
  fileSize: number
  content?: string
}

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
    const allowedTypes = ['application/pdf', 'application/vnd.openxmlformats-officedocument.wordprocessingml.document', 'application/msword', 'text/plain', 'text/markdown']
    const allowedExtensions = ['.pdf', '.docx', '.doc', '.txt', '.md']
    const fileExtension = '.' + file.name.split('.').pop()?.toLowerCase()
    
    if (!allowedTypes.includes(file.type) && !allowedExtensions.includes(fileExtension)) {
      uploadError.value = 'Unsupported file type. Allowed types: PDF, DOCX, DOC, TXT, MD'
      return null
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

