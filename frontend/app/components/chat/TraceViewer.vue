<script setup lang="ts">
interface AgentStep {
  stepIndex: number
  agentName: string
  agentId?: string
  type: string
  action: string
  status: string
  rationale?: string
  observation?: string
  input?: string
  output?: string
  startTime: string
  endTime?: string
  duration?: number
}

interface ErrorInfo {
  type: string
  message: string
  source: string
  stackTrace?: string
}

interface Trace {
  traceId: string
  agentSteps: AgentStep[]
  error?: ErrorInfo
}

interface Props {
  trace: Trace
}

const props = defineProps<Props>()

const isExpanded = ref(true)
const expandedSteps = ref<Set<number>>(new Set())

// Calculate total duration
const totalDuration = computed(() => {
  return props.trace.agentSteps.reduce((sum, step) => sum + (step.duration || 0), 0)
})

const formatTotalDuration = computed(() => {
  const ms = totalDuration.value
  if (ms < 1000) return `${ms}ms`
  if (ms < 60000) return `${(ms / 1000).toFixed(1)}s`
  return `${(ms / 60000).toFixed(1)}m`
})

const toggleStep = (index: number) => {
  if (expandedSteps.value.has(index)) {
    expandedSteps.value.delete(index)
  } else {
    expandedSteps.value.add(index)
  }
  expandedSteps.value = new Set(expandedSteps.value)
}

// Get step display info based on type and action
const getStepDisplay = (step: AgentStep) => {
  const displays: Record<string, { icon: string; label: string; color: string }> = {
    // Thinking steps
    'thinking': {
      icon: '💭',
      label: 'Thinking',
      color: 'text-purple-400'
    },
    // Collaborator calling
    'collaborator_calling': {
      icon: '🤝',
      label: `Calling ${step.agentName}`,
      color: 'text-pink-400'
    },
    // Collaborator response
    'collaborator_response': {
      icon: '📨',
      label: `From ${step.agentName}`,
      color: 'text-pink-400'
    },
    // Action/Tool execution
    'action': {
      icon: '⚡',
      label: step.action,
      color: 'text-yellow-400'
    },
    // Action result
    'action_result': {
      icon: '✅',
      label: 'Action Result',
      color: 'text-green-400'
    },
    // Knowledge base
    'knowledge_base': {
      icon: '🗃️',
      label: step.action,
      color: 'text-cyan-400'
    },
    // Finalizing
    'finalizing': {
      icon: '🎯',
      label: 'Finalizing',
      color: 'text-blue-400'
    },
    // Pre-processing
    'pre_processing': {
      icon: '📥',
      label: 'Processing Input',
      color: 'text-blue-400'
    },
    // Post-processing  
    'post_processing': {
      icon: '📤',
      label: 'Formatting Response',
      color: 'text-green-400'
    },
    // Error
    'error': {
      icon: '❌',
      label: 'Error',
      color: 'text-red-400'
    },
    // Guardrail
    'guardrail': {
      icon: '🛡️',
      label: 'Content Check',
      color: 'text-orange-400'
    },
    // Default
    'default': {
      icon: '⚙️',
      label: step.action || 'Processing',
      color: 'text-gray-400'
    }
  }

  // Determine which display to use
  if (step.type === 'error') return displays.error
  if (step.type === 'guardrail') return displays.guardrail
  if (step.type === 'pre_processing') return displays.pre_processing
  if (step.type === 'post_processing') return displays.post_processing
  if (step.type === 'knowledge_base') return displays.knowledge_base
  
  if (step.type === 'collaborator') {
    if (step.status === 'running' || step.action === 'Calling' || step.action?.includes('Invoking')) {
      return displays.collaborator_calling
    }
    return displays.collaborator_response
  }
  
  if (step.type === 'action') {
    if (step.output) return displays.action_result
    return displays.action
  }
  
  if (step.rationale || step.action === 'Thinking') {
    return displays.thinking
  }
  
  if (step.action === 'Final response') {
    return displays.finalizing
  }

  return displays.default
}

const getStatusBadge = (status: string) => {
  switch (status) {
    case 'success': return { class: 'bg-green-500/20 text-green-400', show: false }
    case 'error': return { class: 'bg-red-500/20 text-red-400', show: true }
    case 'running': return { class: 'bg-blue-500/20 text-blue-400 animate-pulse', show: true }
    default: return { class: 'bg-gray-500/20 text-gray-400', show: false }
  }
}

const formatDuration = (duration?: number) => {
  if (!duration) return ''
  if (duration < 1000) return `${duration}ms`
  return `${(duration / 1000).toFixed(1)}s`
}

const hasDetails = (step: AgentStep) => {
  return step.rationale || step.observation || step.input || step.output
}

const formatActionCall = (step: AgentStep) => {
  if (!step.input) return step.action
  // Try to format as function call
  try {
    const input = JSON.parse(step.input)
    const params = Object.entries(input)
      .map(([k, v]) => `${k}=${typeof v === 'string' ? v.substring(0, 30) + (v.length > 30 ? '...' : '') : v}`)
      .join(', ')
    return `${step.action}(${params})`
  } catch {
    return step.action
  }
}

const formatOutput = (output: string) => {
  try {
    const parsed = JSON.parse(output)
    return JSON.stringify(parsed, null, 2)
  } catch {
    return output
  }
}
</script>

<template>
  <div class="border border-[var(--color-border)] rounded-xl overflow-hidden">
    <!-- Header -->
    <button
      class="w-full px-4 py-2.5 flex items-center justify-between bg-[var(--color-bg-tertiary)] hover:bg-opacity-80 transition-colors"
      @click="isExpanded = !isExpanded"
    >
      <div class="flex items-center gap-2">
        <Icon name="lucide:activity" class="w-4 h-4 text-accent-secondary" />
        <span class="text-sm font-medium">Execution Trace</span>
        <span class="text-xs text-[var(--color-text-muted)] font-mono">{{ trace.traceId }}</span>
      </div>
      <div class="flex items-center gap-3">
        <span 
          v-if="trace.error" 
          class="text-xs px-2 py-0.5 rounded-full bg-red-500/20 text-red-400"
        >
          Error
        </span>
        <span class="text-xs text-[var(--color-text-muted)]">
          {{ trace.agentSteps.length }} steps
        </span>
        <span 
          v-if="totalDuration > 0"
          class="text-xs font-mono px-2 py-0.5 rounded-full bg-accent-secondary/20 text-accent-secondary"
        >
          ⏱️ {{ formatTotalDuration }}
        </span>
        <Icon 
          :name="isExpanded ? 'lucide:chevron-up' : 'lucide:chevron-down'" 
          class="w-4 h-4 text-[var(--color-text-muted)]"
        />
      </div>
    </button>

    <!-- Content -->
    <div v-if="isExpanded" class="bg-[var(--color-bg-secondary)] p-4">
      <!-- Error Display -->
      <div 
        v-if="trace.error" 
        class="mb-4 p-3 rounded-lg bg-red-500/10 border border-red-500/30"
      >
        <div class="flex items-start gap-2">
          <span class="text-lg">❌</span>
          <div class="flex-1 min-w-0">
            <span class="font-medium text-red-400">{{ trace.error.type }}</span>
            <p class="text-sm mt-1 text-[var(--color-text-primary)]">{{ trace.error.message }}</p>
          </div>
        </div>
      </div>

      <!-- Steps Timeline -->
      <div class="space-y-3">
        <div 
          v-for="(step, index) in trace.agentSteps" 
          :key="step.stepIndex"
          class="group"
        >
          <!-- Step Row -->
          <div 
            class="flex items-start gap-3 cursor-pointer hover:bg-[var(--color-bg-tertiary)]/30 rounded-lg p-2 -m-2 transition-colors"
            @click="hasDetails(step) && toggleStep(index)"
          >
            <!-- Icon -->
            <span class="text-lg flex-shrink-0 mt-0.5">{{ getStepDisplay(step)?.icon || '⚙️' }}</span>

            <!-- Content -->
            <div class="flex-1 min-w-0">
              <!-- Agent Name + Action -->
              <div class="flex items-center gap-2 flex-wrap">
                <!-- Agent Name Badge -->
                <span 
                  class="font-semibold text-sm px-2 py-0.5 rounded-md"
                  :class="getStepDisplay(step)?.color || 'text-gray-400'"
                  :style="{ backgroundColor: 'color-mix(in srgb, currentColor 15%, transparent)' }"
                >
                  {{ step.agentName }}
                </span>

                <!-- Arrow -->
                <span class="text-[var(--color-text-muted)]">→</span>

                <!-- Action -->
                <span class="text-sm text-[var(--color-text-secondary)]">
                  {{ step.action || getStepDisplay(step)?.label }}
                </span>
                
                <!-- Status badge (only for running/error) -->
                <span 
                  v-if="getStatusBadge(step.status).show"
                  :class="['text-xs px-1.5 py-0.5 rounded', getStatusBadge(step.status).class]"
                >
                  {{ step.status === 'running' ? 'กำลังทำงาน...' : step.status }}
                </span>

                <!-- Duration -->
                <span 
                  v-if="step.duration" 
                  class="text-xs font-mono px-1.5 py-0.5 rounded bg-[var(--color-bg-tertiary)] text-[var(--color-text-secondary)]"
                >
                  {{ formatDuration(step.duration) }}
                </span>
              </div>

              <!-- Action Call (for action types) -->
              <div 
                v-if="step.type === 'action' && step.action"
                class="mt-1 font-mono text-xs text-[var(--color-text-secondary)] bg-[var(--color-bg-tertiary)] px-2 py-1 rounded inline-block"
              >
                {{ formatActionCall(step) }}
              </div>

              <!-- Rationale/Thinking Preview -->
              <p 
                v-if="step.rationale && !expandedSteps.has(index)" 
                class="text-sm text-[var(--color-text-secondary)] mt-1 line-clamp-2"
              >
                {{ step.rationale }}
              </p>

              <!-- Output Preview (for action results) -->
              <div 
                v-if="step.output && !expandedSteps.has(index) && step.type === 'action'"
                class="mt-1 font-mono text-xs text-[var(--color-text-muted)] bg-[var(--color-bg-tertiary)] px-2 py-1 rounded line-clamp-1"
              >
                {{ step.output.substring(0, 100) }}{{ step.output.length > 100 ? '...' : '' }}
              </div>
            </div>

            <!-- Expand Arrow -->
            <Icon 
              v-if="hasDetails(step)"
              :name="expandedSteps.has(index) ? 'lucide:chevron-up' : 'lucide:chevron-down'" 
              class="w-4 h-4 text-[var(--color-text-muted)] flex-shrink-0 opacity-0 group-hover:opacity-100 transition-opacity"
            />
          </div>

          <!-- Expanded Details -->
          <div 
            v-if="expandedSteps.has(index) && hasDetails(step)" 
            class="ml-9 mt-2 space-y-2"
          >
            <!-- Full Rationale -->
            <div v-if="step.rationale" class="p-3 rounded-lg bg-purple-500/10 border border-purple-500/20">
              <p class="text-sm text-[var(--color-text-primary)] whitespace-pre-wrap">{{ step.rationale }}</p>
            </div>

            <!-- Input -->
            <div v-if="step.input" class="space-y-1">
              <span class="text-xs font-medium text-[var(--color-text-muted)]">Input:</span>
              <pre class="p-2 rounded-lg bg-[var(--color-bg-tertiary)] text-xs overflow-x-auto font-mono text-[var(--color-text-secondary)] whitespace-pre-wrap max-h-40 overflow-y-auto">{{ step.input }}</pre>
            </div>

            <!-- Observation -->
            <div v-if="step.observation" class="p-3 rounded-lg bg-cyan-500/10 border border-cyan-500/20">
              <p class="text-sm text-[var(--color-text-primary)]">{{ step.observation }}</p>
            </div>

            <!-- Output -->
            <div v-if="step.output" class="space-y-1">
              <span class="text-xs font-medium text-green-400">✅ Action Result</span>
              <pre class="p-2 rounded-lg bg-[var(--color-bg-tertiary)] text-xs overflow-x-auto font-mono text-[var(--color-text-secondary)] whitespace-pre-wrap max-h-60 overflow-y-auto">{{ formatOutput(step.output) }}</pre>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div 
        v-if="trace.agentSteps.length === 0 && !trace.error" 
        class="text-center text-[var(--color-text-muted)] py-4"
      >
        <span class="text-2xl">📭</span>
        <p class="text-sm mt-1">No trace data available</p>
      </div>
    </div>
  </div>
</template>
