<script setup lang="ts">
const { sessions, currentSession, createSession, selectSession, deleteSession, isLoading: sessionsLoading } = useSession()
const { messages, isStreaming, thinkingStatus, agentSteps, currentError, wasSummarized, sendMessage, stopStream, clearError, clearHistory } = useChat()

const sidebarOpen = ref(true)
const showClearConfirm = ref(false)

const toggleSidebar = () => {
  sidebarOpen.value = !sidebarOpen.value
}

const handleClearHistory = async () => {
  await clearHistory()
  showClearConfirm.value = false
}

// Auto-scroll to bottom when new messages arrive
const chatContainer = ref<HTMLElement | null>(null)

watch(messages, () => {
  nextTick(() => {
    if (chatContainer.value) {
      chatContainer.value.scrollTop = chatContainer.value.scrollHeight
    }
  })
}, { deep: true })
</script>

<template>
  <div class="h-screen flex overflow-hidden bg-[var(--color-bg-primary)]">
    <!-- Sidebar -->
    <SidebarSessionSidebar
      :open="sidebarOpen"
      :sessions="sessions"
      :current-session-id="currentSession?.id"
      :loading="sessionsLoading"
      @toggle="toggleSidebar"
      @new-chat="createSession"
      @select="selectSession"
      @delete="deleteSession"
    />

    <!-- Main Content -->
    <div class="flex-1 flex flex-col min-w-0">
      <!-- Header -->
      <header class="h-14 flex items-center justify-between px-4 border-b border-[var(--color-border)] bg-[var(--color-bg-secondary)]">
        <div class="flex items-center gap-3">
          <button
            class="btn-icon lg:hidden"
            @click="toggleSidebar"
          >
            <Icon name="lucide:menu" class="w-5 h-5" />
          </button>
          <h1 class="font-semibold text-lg">
            {{ currentSession?.title || 'AgentBedrock UI' }}
          </h1>
        </div>
        <div class="flex items-center gap-2">
          <!-- Clear History Button -->
          <button
            v-if="currentSession && messages.length > 0"
            class="btn-icon text-[var(--color-text-secondary)] hover:text-red-500 transition-colors"
            title="Clear conversation history"
            @click="showClearConfirm = true"
          >
            <Icon name="lucide:trash-2" class="w-5 h-5" />
          </button>
          <UiThemeToggle />
        </div>
      </header>

      <!-- Clear History Confirmation Modal -->
      <Teleport to="body">
        <div 
          v-if="showClearConfirm" 
          class="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
          @click.self="showClearConfirm = false"
        >
          <div class="bg-[var(--color-bg-secondary)] rounded-xl p-6 max-w-md w-full mx-4 shadow-2xl border border-[var(--color-border)]">
            <div class="flex items-center gap-3 mb-4">
              <div class="w-10 h-10 rounded-full bg-red-500/10 flex items-center justify-center">
                <Icon name="lucide:alert-triangle" class="w-5 h-5 text-red-500" />
              </div>
              <h3 class="text-lg font-semibold">Clear Conversation History</h3>
            </div>
            <p class="text-[var(--color-text-secondary)] mb-6">
              This will delete all messages in this conversation. This action cannot be undone.
            </p>
            <div class="flex gap-3 justify-end">
              <button 
                class="px-4 py-2 rounded-lg border border-[var(--color-border)] hover:bg-[var(--color-bg-tertiary)] transition-colors"
                @click="showClearConfirm = false"
              >
                Cancel
              </button>
              <button 
                class="px-4 py-2 rounded-lg bg-red-500 text-white hover:bg-red-600 transition-colors"
                @click="handleClearHistory"
              >
                Clear History
              </button>
            </div>
          </div>
        </div>
      </Teleport>

      <!-- Chat Area -->
      <main 
        ref="chatContainer"
        class="flex-1 overflow-y-auto scrollbar-thin"
      >
        <div class="max-w-4xl mx-auto py-6 px-4">
          <!-- Empty State -->
          <div v-if="!currentSession" class="flex flex-col items-center justify-center h-[60vh] text-center">
            <div class="w-20 h-20 rounded-2xl bg-gradient-to-br from-accent-primary to-accent-secondary flex items-center justify-center mb-6">
              <Icon name="lucide:bot" class="w-10 h-10 text-white" />
            </div>
            <h2 class="text-2xl font-semibold mb-2">Welcome to AgentBedrock UI</h2>
            <p class="text-[var(--color-text-secondary)] mb-6 max-w-md">
              An easy-to-use interface for testing AWS AgentBedrock Team Agents
            </p>
            <button class="btn-primary" @click="createSession()">
              <Icon name="lucide:plus" class="w-4 h-4 mr-2" />
              Start New Chat
            </button>
          </div>

          <!-- Messages -->
          <ChatMessageList
            v-else
            :messages="messages"
            :is-streaming="isStreaming"
            :thinking-status="thinkingStatus"
            :agent-steps="agentSteps"
            :error="currentError"
            :was-summarized="wasSummarized"
            @dismiss-error="clearError"
            @dismiss-summarized="wasSummarized = false"
          />
        </div>
      </main>

      <!-- Input Area -->
      <div v-if="currentSession" class="border-t border-[var(--color-border)] bg-[var(--color-bg-secondary)]">
        <div class="max-w-4xl mx-auto py-4 px-4">
          <ChatInput
            :disabled="!currentSession"
            :is-streaming="isStreaming"
            @send="sendMessage"
            @stop="stopStream"
          />
        </div>
      </div>
    </div>
  </div>
</template>

