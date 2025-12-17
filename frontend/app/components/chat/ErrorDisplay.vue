<script setup lang="ts">
interface ErrorInfo {
  type: string
  message: string
  source?: string
  stackTrace?: string
}

interface Props {
  error: ErrorInfo
}

const props = defineProps<Props>()
const emit = defineEmits<{
  dismiss: []
}>()

const copied = ref(false)

const errorText = computed(() => {
  let text = `Error Type: ${props.error.type}\n`
  text += `Message: ${props.error.message}\n`
  if (props.error.source) {
    text += `Source: ${props.error.source}\n`
  }
  if (props.error.stackTrace) {
    text += `\nStack Trace:\n${props.error.stackTrace}`
  }
  return text
})

const copyError = async () => {
  try {
    await navigator.clipboard.writeText(errorText.value)
    copied.value = true
    setTimeout(() => {
      copied.value = false
    }, 2000)
  } catch (err) {
    console.error('Failed to copy:', err)
  }
}
</script>

<template>
  <div class="rounded-xl border border-accent-error/30 bg-accent-error/10 overflow-hidden animate-fade-in">
    <!-- Header -->
    <div class="flex items-center justify-between px-4 py-2 bg-accent-error/20 border-b border-accent-error/30">
      <div class="flex items-center gap-2">
        <Icon name="lucide:alert-triangle" class="w-5 h-5 text-accent-error" />
        <span class="font-medium text-accent-error">{{ error.type }}</span>
      </div>
      <div class="flex items-center gap-2">
        <button
          @click="copyError"
          class="flex items-center gap-1.5 px-2 py-1 rounded-lg text-sm transition-all"
          :class="copied 
            ? 'bg-accent-primary/20 text-accent-primary' 
            : 'bg-[var(--color-bg-tertiary)] hover:bg-[var(--color-border)] text-[var(--color-text-secondary)]'"
        >
          <Icon :name="copied ? 'lucide:check' : 'lucide:copy'" class="w-4 h-4" />
          <span>{{ copied ? 'Copied!' : 'Copy' }}</span>
        </button>
        <button
          @click="emit('dismiss')"
          class="flex items-center justify-center w-7 h-7 rounded-lg text-sm transition-all bg-[var(--color-bg-tertiary)] hover:bg-[var(--color-border)] text-[var(--color-text-secondary)]"
          title="Dismiss"
        >
          <Icon name="lucide:x" class="w-4 h-4" />
        </button>
      </div>
    </div>

    <!-- Error Content -->
    <div class="p-4 space-y-3">
      <!-- Message -->
      <div>
        <p class="text-sm text-[var(--color-text-primary)]">{{ error.message }}</p>
      </div>

      <!-- Source -->
      <div v-if="error.source" class="flex items-center gap-2 text-sm">
        <span class="text-[var(--color-text-muted)]">Source:</span>
        <code class="px-2 py-0.5 rounded bg-[var(--color-bg-tertiary)] text-accent-secondary font-mono text-xs">
          {{ error.source }}
        </code>
      </div>

      <!-- Stack Trace -->
      <div v-if="error.stackTrace" class="space-y-1">
        <span class="text-xs text-[var(--color-text-muted)]">Stack Trace:</span>
        <pre class="p-3 rounded-lg bg-dark-900 text-dark-200 text-xs overflow-x-auto font-mono scrollbar-thin">{{ error.stackTrace }}</pre>
      </div>

      <!-- Quick Actions -->
      <div class="flex items-center gap-2 pt-2 border-t border-accent-error/20">
        <span class="text-xs text-[var(--color-text-muted)]">💡 Tip:</span>
        <span class="text-xs text-[var(--color-text-secondary)]">
          Check AWS credentials, Agent ID, and Agent Alias configuration
        </span>
      </div>
    </div>
  </div>
</template>

