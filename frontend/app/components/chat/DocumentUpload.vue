<script setup lang="ts">
interface Props {
  sessionId: string
  disabled?: boolean
}

const props = defineProps<Props>()

const { uploadFile, isUploading, uploadError, uploadProgress } = useDocumentUpload()
const { currentSession } = useSession()

const fileInputRef = ref<HTMLInputElement | null>(null)
const isDragging = ref(false)

const handleFileSelect = async (files: FileList | null) => {
  if (!files || files.length === 0) return
  
  const sessionId = props.sessionId || currentSession.value?.id
  if (!sessionId) {
    alert('Please select a session first')
    return
  }

  // Upload each file
  for (let i = 0; i < files.length; i++) {
    await uploadFile(files[i], sessionId)
  }
}

const handleFileInput = (e: Event) => {
  const target = e.target as HTMLInputElement
  handleFileSelect(target.files)
  // Reset input
  if (target) {
    target.value = ''
  }
}

const handleDrop = (e: DragEvent) => {
  e.preventDefault()
  isDragging.value = false
  
  if (props.disabled) return
  
  const files = e.dataTransfer?.files
  handleFileSelect(files)
}

const handleDragOver = (e: DragEvent) => {
  e.preventDefault()
  if (!props.disabled) {
    isDragging.value = true
  }
}

const handleDragLeave = () => {
  isDragging.value = false
}

const openFilePicker = () => {
  fileInputRef.value?.click()
}

const getFileTypeIcon = (fileType: string) => {
  switch (fileType) {
    case 'pdf':
      return 'lucide:file-text'
    case 'docx':
    case 'doc':
      return 'lucide:file-word'
    case 'txt':
    case 'md':
      return 'lucide:file-code'
    default:
      return 'lucide:file'
  }
}

watch(uploadError, (error) => {
  if (error) {
    console.error('Upload error:', error)
  }
})
</script>

<template>
  <div class="relative">
    <input
      ref="fileInputRef"
      type="file"
      accept=".pdf,.docx,.doc,.txt,.md"
      multiple
      class="hidden"
      :disabled="disabled"
      @change="handleFileInput"
    />
    
    <div
      :class="[
        'border-2 border-dashed rounded-xl p-4 transition-all cursor-pointer',
        isDragging
          ? 'border-accent-primary bg-accent-primary/10'
          : 'border-[var(--color-border)] hover:border-accent-primary/50',
        disabled && 'opacity-50 cursor-not-allowed'
      ]"
      @click="openFilePicker"
      @drop="handleDrop"
      @dragover="handleDragOver"
      @dragleave="handleDragLeave"
    >
      <div class="flex flex-col items-center justify-center gap-2 text-center">
        <Icon
          name="lucide:upload"
          class="w-8 h-8 text-[var(--color-text-muted)]"
        />
        <div class="text-sm">
          <span class="text-[var(--color-text-primary)] font-medium">
            Click to upload
          </span>
          <span class="text-[var(--color-text-muted)]"> or drag and drop</span>
        </div>
        <p class="text-xs text-[var(--color-text-muted)]">
          PDF, DOCX, DOC, TXT, MD (max 10MB)
        </p>
      </div>
    </div>

    <!-- Upload Progress -->
    <div v-if="isUploading" class="mt-2">
      <div class="flex items-center gap-2">
        <div class="flex-1 h-2 bg-[var(--color-bg-tertiary)] rounded-full overflow-hidden">
          <div
            class="h-full bg-accent-primary transition-all duration-300"
            :style="{ width: `${uploadProgress}%` }"
          />
        </div>
        <span class="text-xs text-[var(--color-text-muted)]">{{ uploadProgress }}%</span>
      </div>
    </div>

    <!-- Upload Error -->
    <div v-if="uploadError" class="mt-2 p-2 bg-accent-error/10 border border-accent-error/20 rounded-lg">
      <p class="text-sm text-accent-error">{{ uploadError }}</p>
    </div>
  </div>
</template>

