<script setup lang="ts">
interface Props {
  disabled?: boolean
  isStreaming?: boolean
  sessionId?: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  send: [message: string]
  stop: []
}>()

// Get document upload composable
const documentUpload = useDocumentUpload()
const uploadedDocuments = documentUpload.uploadedDocuments
const removeDocument = documentUpload.removeDocument
const clearDocuments = documentUpload.clearDocuments
const showUpload = ref(false)

const toggleUpload = () => {
  showUpload.value = !showUpload.value
}

const message = ref('')
const textareaRef = ref<HTMLTextAreaElement | null>(null)

// Auto resize textarea
const autoResize = () => {
  if (textareaRef.value) {
    textareaRef.value.style.height = 'auto'
    textareaRef.value.style.height = Math.min(textareaRef.value.scrollHeight, 200) + 'px'
  }
}

const handleSubmit = () => {
  if (props.isStreaming) {
    emit('stop')
    return
  }
  
  // Allow sending if there's a message OR documents attached
  if ((!message.value.trim() && uploadedDocuments.value.length === 0) || props.disabled) {
    return
  }
  
  emit('send', message.value || '') // Send empty string if only documents
  message.value = ''
  clearDocuments()
  showUpload.value = false
  nextTick(autoResize)
}

const handleKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    handleSubmit()
  }
}

watch(message, () => {
  nextTick(autoResize)
})
</script>

<template>
  <div class="relative">
    <!-- Document Upload Area & Document List -->
    <ClientOnly>
      <div v-show="showUpload" aria-label="Document Upload Area" class="mb-3 p-4 bg-[var(--color-bg-tertiary)] border border-dashed border-[var(--color-border)] rounded-xl">
        <ChatDocumentUpload :session-id="sessionId || ''" :disabled="disabled" />
      </div>
      
      <!-- Document List -->
      <div v-if="uploadedDocuments && uploadedDocuments.length > 0" class="mb-2">
        <ChatDocumentList :documents="uploadedDocuments" @remove="removeDocument" />
      </div>
    </ClientOnly>

    <div class="relative flex items-end gap-2 p-2 rounded-2xl bg-[var(--color-bg-tertiary)] border border-[var(--color-border)] focus-within:ring-2 focus-within:ring-accent-primary/50 focus-within:border-accent-primary transition-all">
      <!-- Paperclip Button -->
      <button
        type="button"
        :disabled="disabled"
        :class="[
          'flex-shrink-0 w-10 h-10 rounded-xl flex items-center justify-center transition-colors disabled:opacity-50 disabled:cursor-not-allowed',
          showUpload 
            ? 'bg-accent-primary/20 text-accent-primary' 
            : 'bg-[var(--color-bg-secondary)] hover:bg-[var(--color-bg-secondary)]/80 text-[var(--color-text-muted)]'
        ]"
        @click="toggleUpload"
      >
        <Icon name="lucide:paperclip" class="w-5 h-5" />
      </button>

      <!-- Textarea -->
      <textarea
        ref="textareaRef"
        v-model="message"
        :disabled="disabled"
        placeholder="Send a message..."
        rows="1"
        class="flex-1 resize-none bg-transparent border-0 focus:outline-none focus:ring-0 text-[var(--color-text-primary)] placeholder:text-[var(--color-text-muted)] py-2 px-2 max-h-[200px] scrollbar-thin disabled:opacity-50 disabled:cursor-not-allowed"
        @keydown="handleKeydown"
      />
      
      <!-- Send/Stop Button -->
      <button
        type="button"
        :disabled="disabled || (!isStreaming && !message.trim() && (!uploadedDocuments || uploadedDocuments.length === 0))"
        :class="[
          'flex-shrink-0 w-10 h-10 rounded-xl flex items-center justify-center transition-all duration-150 disabled:opacity-50 disabled:cursor-not-allowed',
          isStreaming
            ? 'bg-red-500 hover:bg-red-600 text-white'
            : 'bg-accent-primary hover:bg-accent-primary/90 text-white'
        ]"
        @click="handleSubmit"
      >
        <Icon v-if="isStreaming" name="lucide:square" class="w-4 h-4" />
        <Icon v-else name="lucide:arrow-up" class="w-5 h-5" />
      </button>
    </div>
    
    <p class="text-xs text-[var(--color-text-muted)] text-center mt-2">
      Press <kbd class="px-1.5 py-0.5 rounded bg-[var(--color-bg-tertiary)] text-xs font-mono">Enter</kbd> to send, 
      <kbd class="px-1.5 py-0.5 rounded bg-[var(--color-bg-tertiary)] text-xs font-mono">Shift+Enter</kbd> for new line
    </p>
  </div>
</template>
