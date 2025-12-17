<script setup lang="ts">
interface Session {
  id: string
  title: string
  createdAt: string
  updatedAt: string
}

interface Props {
  open: boolean
  sessions: Session[]
  currentSessionId?: string
  loading?: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
  toggle: []
  newChat: []
  select: [id: string]
  delete: [id: string]
}>()

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  
  if (days === 0) return 'Today'
  if (days === 1) return 'Yesterday'
  if (days < 7) return `${days} days ago`
  return date.toLocaleDateString()
}

const confirmDelete = (id: string, event: Event) => {
  event.stopPropagation()
  if (confirm('Delete this conversation?')) {
    emit('delete', id)
  }
}
</script>

<template>
  <aside
    :class="[
      'fixed lg:relative inset-y-0 left-0 z-40 w-72 flex flex-col',
      'bg-[var(--color-bg-secondary)] border-r border-[var(--color-border)]',
      'transition-transform duration-300 ease-in-out',
      open ? 'translate-x-0' : '-translate-x-full lg:translate-x-0 lg:w-0 lg:opacity-0 lg:invisible'
    ]"
  >
    <!-- Header -->
    <div class="h-14 flex items-center justify-between px-4 border-b border-[var(--color-border)]">
      <div class="flex items-center gap-2">
        <div class="w-8 h-8 rounded-lg bg-gradient-to-br from-accent-primary to-accent-secondary flex items-center justify-center">
          <Icon name="lucide:bot" class="w-4 h-4 text-white" />
        </div>
        <span class="font-semibold text-sm">AgentBedrock</span>
      </div>
      <button class="btn-icon" @click="$emit('toggle')">
        <Icon name="lucide:panel-left-close" class="w-5 h-5" />
      </button>
    </div>

    <!-- New Chat Button -->
    <div class="p-3">
      <button 
        class="w-full btn-secondary flex items-center justify-center gap-2"
        @click="$emit('newChat')"
      >
        <Icon name="lucide:plus" class="w-4 h-4" />
        New Chat
      </button>
    </div>

    <!-- Sessions List -->
    <div class="flex-1 overflow-y-auto scrollbar-thin px-2 pb-4">
      <!-- Loading -->
      <div v-if="loading" class="flex items-center justify-center py-8">
        <div class="flex gap-1">
          <div class="typing-dot"></div>
          <div class="typing-dot"></div>
          <div class="typing-dot"></div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else-if="sessions.length === 0" class="text-center py-8 px-4">
        <Icon name="lucide:message-square" class="w-12 h-12 mx-auto mb-3 text-[var(--color-text-muted)]" />
        <p class="text-sm text-[var(--color-text-secondary)]">No conversations yet</p>
        <p class="text-xs text-[var(--color-text-muted)] mt-1">Start a new chat to begin</p>
      </div>

      <!-- Sessions -->
      <div v-else class="space-y-1">
        <button
          v-for="session in sessions"
          :key="session.id"
          :class="[
            'w-full text-left px-3 py-2.5 rounded-xl transition-all duration-150 group',
            'hover:bg-[var(--color-bg-tertiary)]',
            session.id === currentSessionId
              ? 'bg-[var(--color-bg-tertiary)] ring-1 ring-accent-primary/30'
              : ''
          ]"
          @click="$emit('select', session.id)"
        >
          <div class="flex items-start justify-between gap-2">
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium truncate">{{ session.title }}</p>
              <p class="text-xs text-[var(--color-text-muted)] mt-0.5">
                {{ formatDate(session.updatedAt) }}
              </p>
            </div>
            <button
              class="opacity-0 group-hover:opacity-100 p-1 rounded hover:bg-accent-error/10 hover:text-accent-error transition-all"
              @click="confirmDelete(session.id, $event)"
            >
              <Icon name="lucide:trash-2" class="w-4 h-4" />
            </button>
          </div>
        </button>
      </div>
    </div>

    <!-- Footer -->
    <div class="p-3 border-t border-[var(--color-border)]">
      <div class="text-xs text-[var(--color-text-muted)] text-center">
        AgentBedrock UI v1.0
      </div>
    </div>
  </aside>

  <!-- Overlay for mobile -->
  <div
    v-if="open"
    class="fixed inset-0 bg-black/50 z-30 lg:hidden"
    @click="$emit('toggle')"
  />
</template>

