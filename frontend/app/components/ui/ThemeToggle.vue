<script setup lang="ts">
const { rawTheme, toggleTheme } = useTheme()

const iconName = computed(() => {
  switch (rawTheme.value) {
    case 'light': return 'lucide:sun'
    case 'dark': return 'lucide:moon'
    default: return 'lucide:monitor'
  }
})

const label = computed(() => {
  switch (rawTheme.value) {
    case 'light': return 'Light'
    case 'dark': return 'Dark'
    default: return 'System'
  }
})
</script>

<template>
  <button 
    class="btn-icon flex items-center gap-2 text-sm"
    @click="toggleTheme"
    :title="`Theme: ${label}`"
  >
    <Icon :name="iconName" class="w-5 h-5" />
    <span class="hidden sm:inline text-[var(--color-text-secondary)]">{{ label }}</span>
  </button>
</template>

