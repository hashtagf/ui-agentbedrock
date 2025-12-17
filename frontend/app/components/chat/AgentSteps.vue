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
      v-for="step in steps" 
      :key="step.stepIndex"
      class="animate-fade-in"
    >
      <!-- Step Row -->
      <div class="flex items-start gap-2">
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

          <!-- Rationale Preview -->
          <p 
            v-if="step.rationale" 
            class="text-sm text-[var(--color-text-secondary)] mt-1 line-clamp-2"
          >
            {{ step.rationale }}
          </p>

          <!-- Action Call -->
          <div 
            v-if="step.type === 'action' && step.action && step.input"
            class="mt-1 font-mono text-xs text-[var(--color-text-muted)] bg-[var(--color-bg-tertiary)] px-2 py-1 rounded inline-block"
          >
            {{ step.action }}(...)
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
