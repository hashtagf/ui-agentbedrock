<script setup lang="ts">
interface Document {
  documentId: string
  filename: string
  fileType: string
  fileSize: number
  content?: string
}

interface Props {
  documents: Document[]
  showRemove?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  showRemove: true,
})

const emit = defineEmits<{
  remove: [documentId: string]
}>()

const formatFileSize = (bytes: number): string => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
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

const getFileTypeColor = (fileType: string) => {
  switch (fileType) {
    case 'pdf':
      return 'text-red-500'
    case 'docx':
    case 'doc':
      return 'text-blue-500'
    case 'txt':
    case 'md':
      return 'text-green-500'
    default:
      return 'text-[var(--color-text-muted)]'
  }
}

const handleRemove = (documentId: string) => {
  emit('remove', documentId)
}
</script>

<template>
  <div v-if="documents.length > 0" class="flex flex-wrap gap-2">
    <div
      v-for="doc in documents"
      :key="doc.documentId"
      class="flex items-center gap-2 px-3 py-2 bg-[var(--color-bg-tertiary)] rounded-lg border border-[var(--color-border)]"
    >
      <Icon
        :name="getFileTypeIcon(doc.fileType)"
        :class="['w-4 h-4', getFileTypeColor(doc.fileType)]"
      />
      <div class="flex-1 min-w-0">
        <p class="text-sm font-medium text-[var(--color-text-primary)] truncate">
          {{ doc.filename }}
        </p>
        <p class="text-xs text-[var(--color-text-muted)]">
          {{ formatFileSize(doc.fileSize) }}
        </p>
      </div>
      <button
        v-if="showRemove"
        @click="handleRemove(doc.documentId)"
        class="flex-shrink-0 p-1 hover:bg-[var(--color-bg-secondary)] rounded transition-colors"
        title="Remove document"
      >
        <Icon name="lucide:x" class="w-4 h-4 text-[var(--color-text-muted)]" />
      </button>
    </div>
  </div>
</template>

