<script setup lang="ts">
import MarkdownIt from 'markdown-it'

interface Message {
  id: string
  role: 'user' | 'assistant'
  content: string
  trace?: any
  createdAt: string
}

interface AgentStep {
  stepIndex: number
  agentName: string
  action: string
  status: 'running' | 'success' | 'error'
}

interface ErrorInfo {
  type: string
  message: string
  source?: string
  stackTrace?: string
}

interface Props {
  messages: Message[]
  isStreaming?: boolean
  thinkingStatus?: string | null
  agentSteps?: AgentStep[]
  error?: ErrorInfo | null
  wasSummarized?: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
  dismissError: []
  dismissSummarized: []
}>()

const md = new MarkdownIt({
  html: false,
  linkify: true,
  typographer: true,
  breaks: true,
})

const renderMarkdown = (content: string) => {
  return md.render(content)
}

const formatTime = (dateStr: string) => {
  return new Date(dateStr).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}
</script>

<template>
  <div class="space-y-6">
    <!-- Summarization Notice -->
    <div v-if="wasSummarized" class="animate-fade-in">
      <div class="bg-blue-500/10 border border-blue-500/30 rounded-xl p-4">
        <div class="flex items-start gap-3">
          <div class="flex-shrink-0 w-8 h-8 rounded-lg bg-blue-500/20 flex items-center justify-center">
            <Icon name="lucide:sparkles" class="w-4 h-4 text-blue-400" />
          </div>
          <div class="flex-1">
            <h4 class="font-medium text-blue-400 mb-1">Context Window Optimized</h4>
            <p class="text-sm text-[var(--color-text-secondary)]">
              The conversation was automatically summarized to avoid token limits. 
              A new agent session was created with the summarized context, so you can continue seamlessly.
            </p>
          </div>
          <button 
            class="text-[var(--color-text-muted)] hover:text-[var(--color-text-primary)] transition-colors"
            @click="emit('dismissSummarized')"
          >
            <Icon name="lucide:x" class="w-4 h-4" />
          </button>
        </div>
      </div>
    </div>

    <!-- Messages -->
    <div
      v-for="message in messages"
      :key="message.id"
      class="animate-fade-in"
    >
      <!-- User Message -->
      <div v-if="message.role === 'user'" class="flex justify-end">
        <div class="max-w-[85%] lg:max-w-[70%]">
          <div class="message-user rounded-2xl rounded-tr-md px-4 py-3">
            <p class="text-[var(--color-text-primary)] whitespace-pre-wrap">{{ message.content }}</p>
          </div>
          <p class="text-xs text-[var(--color-text-muted)] mt-1 text-right">
            {{ formatTime(message.createdAt) }}
          </p>
        </div>
      </div>

      <!-- Assistant Message -->
      <div v-else class="flex justify-start">
        <div class="max-w-[85%] lg:max-w-[70%]">
          <div class="flex items-start gap-3">
            <!-- Avatar -->
            <div class="flex-shrink-0 w-8 h-8 rounded-lg bg-gradient-to-br from-accent-primary to-accent-secondary flex items-center justify-center">
              <Icon name="lucide:bot" class="w-4 h-4 text-white" />
            </div>
            
            <div class="flex-1 min-w-0">
              <!-- Thinking Status -->
              <div v-if="isStreaming && message.id.startsWith('temp-assistant') && thinkingStatus" class="mb-3 animate-fade-in">
                <div class="flex items-center gap-2 text-sm text-[var(--color-text-secondary)]">
                  <div class="flex gap-1">
                    <div class="typing-dot"></div>
                    <div class="typing-dot"></div>
                    <div class="typing-dot"></div>
                  </div>
                  <span class="italic">{{ thinkingStatus }}</span>
                </div>
              </div>

              <!-- Agent Steps (realtime during streaming) -->
              <ChatAgentSteps 
                v-if="isStreaming && message.id.startsWith('temp-assistant') && agentSteps && agentSteps.length > 0"
                :steps="agentSteps"
                class="mb-3"
              />

              <!-- Trace Viewer (ABOVE message content) -->
              <ChatTraceViewer 
                v-if="message.trace" 
                :trace="message.trace"
                class="mb-3"
              />

              <!-- Message Content -->
              <div 
                v-if="message.content"
                class="message-assistant rounded-2xl rounded-tl-md px-4 py-3"
              >
                <div 
                  class="prose-chat"
                  v-html="renderMarkdown(message.content)"
                />
              </div>

              <!-- Loading placeholder when empty -->
              <div 
                v-else-if="isStreaming && message.id.startsWith('temp-assistant')"
                class="message-assistant rounded-2xl rounded-tl-md px-4 py-3"
              >
                <div class="flex gap-1">
                  <div class="typing-dot"></div>
                  <div class="typing-dot"></div>
                  <div class="typing-dot"></div>
                </div>
              </div>

              <p v-if="message.content" class="text-xs text-[var(--color-text-muted)] mt-1">
                {{ formatTime(message.createdAt) }}
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Error Display -->
    <div v-if="error" class="animate-fade-in">
      <div class="flex justify-start">
        <div class="max-w-[90%] lg:max-w-[80%]">
          <div class="flex items-start gap-3">
            <!-- Avatar -->
            <div class="flex-shrink-0 w-8 h-8 rounded-lg bg-gradient-to-br from-accent-error to-red-700 flex items-center justify-center">
              <Icon name="lucide:alert-triangle" class="w-4 h-4 text-white" />
            </div>
            
            <div class="flex-1 min-w-0">
              <ChatErrorDisplay :error="error" @dismiss="emit('dismissError')" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

