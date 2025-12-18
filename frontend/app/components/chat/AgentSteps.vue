<script setup lang="ts">
interface AgentStep {
  stepIndex: number
  agentName: string
  agentId?: string
  type?: string
  action: string
  status: 'running' | 'success' | 'error'
  rationale?: string
  observation?: string
  input?: string
  output?: string
  duration?: number
}

interface Props {
  steps: AgentStep[]
}

defineProps<Props>()

const expandedSteps = ref<Set<number>>(new Set())

const toggleStep = (index: number) => {
  if (expandedSteps.value.has(index)) {
    expandedSteps.value.delete(index)
  } else {
    expandedSteps.value.add(index)
  }
  expandedSteps.value = new Set(expandedSteps.value)
}

const hasDetails = (step: AgentStep) => {
  return step.rationale || step.observation || step.input || step.output
}

const getStepDisplay = (step: AgentStep) => {
  // Default display
  let icon = '⚙️'
  let label = step.action || 'Processing'
  let color = 'text-gray-400'
  let sublabel = ''

  if (step.type === 'error') {
    icon = '❌'
    label = 'Error'
    color = 'text-red-400'
  } else if (step.type === 'guardrail') {
    icon = '🛡️'
    label = 'Content Check'
    color = 'text-orange-400'
  } else if (step.type === 'pre_processing') {
    icon = '📥'
    label = 'Processing Input'
    color = 'text-blue-400'
  } else if (step.type === 'post_processing') {
    icon = '📤'
    label = 'Formatting Response'
    color = 'text-green-400'
  } else if (step.type === 'knowledge_base') {
    icon = '🗃️'
    label = step.action || 'Knowledge Base'
    color = 'text-cyan-400'
  } else if (step.type === 'collaborator') {
    if (step.status === 'running') {
      icon = '🤝'
      label = `Calling ${step.agentName}`
      color = 'text-pink-400'
      sublabel = 'กำลังเรียกใช้ agent...'
    } else {
      icon = '📨'
      label = `From ${step.agentName}`
      color = 'text-pink-400'
    }
  } else if (step.type === 'action') {
    if (step.status === 'running') {
      icon = '⚡'
      label = step.action || 'Executing action'
      color = 'text-yellow-400'
    } else {
      icon = '✅'
      label = 'Action Result'
      color = 'text-green-400'
    }
  } else if (step.rationale || step.action === 'Thinking') {
    icon = '💭'
    label = 'Thinking'
    color = 'text-purple-400'
  } else if (step.action === 'Final response') {
    icon = '🎯'
    label = 'Finalizing'
    color = 'text-blue-400'
    sublabel = 'กำลังสร้างคำตอบ...'
  } else if (step.action === 'Model invocation') {
    icon = '🧠'
    label = 'Processing'
    color = 'text-purple-400'
  }

  return { icon, label, color, sublabel }
}

const formatDuration = (duration?: number) => {
  if (!duration) return ''
  if (duration < 1000) return `${duration}ms`
  return `${(duration / 1000).toFixed(1)}s`
}
</script>

<template>
  <div class="space-y-2 py-2">
    <div 
      v-for="(step, index) in steps" 
      :key="step.stepIndex"
      class="animate-fade-in group"
    >
      <!-- Step Row -->
      <div 
        class="flex items-start gap-2 cursor-pointer hover:bg-[var(--color-bg-tertiary)]/30 rounded-lg p-2 -m-2 transition-colors"
        :class="{ 'cursor-pointer': hasDetails(step) }"
        @click="hasDetails(step) && toggleStep(index)"
      >
        <!-- Icon -->
        <span class="text-base flex-shrink-0">{{ getStepDisplay(step).icon }}</span>

        <!-- Content -->
        <div class="flex-1 min-w-0">
          <div class="flex items-center gap-2 flex-wrap">
            <!-- Agent Name Badge -->
            <span 
              class="font-semibold text-sm px-2 py-0.5 rounded-md"
              :class="getStepDisplay(step).color"
              :style="{ backgroundColor: 'color-mix(in srgb, currentColor 15%, transparent)' }"
            >
              {{ step.agentName }}
            </span>

            <!-- Arrow -->
            <span class="text-[var(--color-text-muted)]">→</span>

            <!-- Action -->
            <span class="text-sm text-[var(--color-text-secondary)]">
              {{ step.action || getStepDisplay(step).label }}
            </span>
            
            <!-- Running indicator -->
            <span 
              v-if="step.status === 'running'" 
              class="text-xs px-1.5 py-0.5 rounded bg-blue-500/20 text-blue-400 animate-pulse"
            >
              กำลังทำงาน...
            </span>

            <!-- Duration -->
            <span 
              v-if="step.duration && step.status !== 'running'" 
              class="text-xs font-mono px-1.5 py-0.5 rounded bg-[var(--color-bg-tertiary)] text-[var(--color-text-secondary)]"
            >
              {{ formatDuration(step.duration) }}
            </span>
          </div>

          <!-- Sub-label -->
          <p 
            v-if="getStepDisplay(step).sublabel && step.status === 'running'" 
            class="text-xs text-[var(--color-text-muted)] mt-0.5"
          >
            {{ getStepDisplay(step).sublabel }}
          </p>

          <!-- Rationale Preview (collapsed) -->
          <p 
            v-if="step.rationale && !expandedSteps.has(index)" 
            class="text-sm text-[var(--color-text-secondary)] mt-1 line-clamp-2"
          >
            {{ step.rationale }}
          </p>

          <!-- Input Preview (for collaborator calls - collapsed) -->
          <div 
            v-if="step.type === 'collaborator' && step.input && !expandedSteps.has(index)" 
            class="mt-1 text-xs text-[var(--color-text-muted)] bg-[var(--color-bg-tertiary)] px-2 py-1 rounded line-clamp-2"
          >
            <span class="font-medium">📤 Request:</span> {{ step.input.length > 150 ? step.input.substring(0, 150) + '...' : step.input }}
          </div>

          <!-- Output Preview (for collaborator responses - collapsed) -->
          <div 
            v-if="step.type === 'collaborator' && step.output && !expandedSteps.has(index)" 
            class="mt-1 text-xs text-[var(--color-text-muted)] bg-[var(--color-bg-tertiary)] px-2 py-1 rounded line-clamp-2"
          >
            <span class="font-medium">📥 Response:</span> {{ step.output.length > 150 ? step.output.substring(0, 150) + '...' : step.output }}
          </div>

          <!-- Action Call -->
          <div 
            v-if="step.type === 'action' && step.action && step.input"
            class="mt-1 font-mono text-xs text-[var(--color-text-muted)] bg-[var(--color-bg-tertiary)] px-2 py-1 rounded inline-block"
          >
            {{ step.action }}(...)
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
        class="ml-6 mt-2 space-y-2"
      >
        <!-- Full Rationale -->
        <div v-if="step.rationale" class="p-3 rounded-lg bg-purple-500/10 border border-purple-500/20">
          <p class="text-xs font-medium text-purple-400 mb-1">💭 Thinking:</p>
          <p class="text-sm text-[var(--color-text-primary)] whitespace-pre-wrap">{{ step.rationale }}</p>
        </div>

        <!-- Input (for collaborator calls) -->
        <div v-if="step.input" class="space-y-1">
          <span class="text-xs font-medium text-[var(--color-text-muted)]">📤 Input:</span>
          <pre class="p-2 rounded-lg bg-[var(--color-bg-tertiary)] text-xs overflow-x-auto font-mono text-[var(--color-text-secondary)] whitespace-pre-wrap max-h-40 overflow-y-auto">{{ step.input }}</pre>
        </div>

        <!-- Observation -->
        <div v-if="step.observation" class="p-3 rounded-lg bg-cyan-500/10 border border-cyan-500/20">
          <p class="text-xs font-medium text-cyan-400 mb-1">👁️ Observation:</p>
          <p class="text-sm text-[var(--color-text-primary)]">{{ step.observation }}</p>
        </div>

        <!-- Output (for collaborator responses) -->
        <div v-if="step.output" class="space-y-1">
          <span class="text-xs font-medium text-green-400">📥 Output:</span>
          <pre class="p-2 rounded-lg bg-[var(--color-bg-tertiary)] text-xs overflow-x-auto font-mono text-[var(--color-text-secondary)] whitespace-pre-wrap max-h-60 overflow-y-auto">{{ step.output }}</pre>
        </div>
      </div>
    </div>
  </div>
</template>
