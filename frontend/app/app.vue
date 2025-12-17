<script setup lang="ts">
const { theme } = useTheme()
</script>

<template>
  <div :class="theme">
    <NuxtPage />
  </div>
</template>

