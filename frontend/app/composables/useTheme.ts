type Theme = 'light' | 'dark' | 'system'

export function useTheme() {
  const theme = useState<Theme>('theme', () => 'system')
  const resolvedTheme = useState<'light' | 'dark'>('resolvedTheme', () => 'dark')

  const updateResolvedTheme = () => {
    if (import.meta.client) {
      if (theme.value === 'system') {
        resolvedTheme.value = window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
      } else {
        resolvedTheme.value = theme.value
      }
      
      // Update class on documentElement
      if (resolvedTheme.value === 'dark') {
        document.documentElement.classList.add('dark')
      } else {
        document.documentElement.classList.remove('dark')
      }
      
      // Persist to localStorage
      localStorage.setItem('theme', theme.value)
    }
  }

  const setTheme = (newTheme: Theme) => {
    theme.value = newTheme
    updateResolvedTheme()
  }

  const toggleTheme = () => {
    const themes: Theme[] = ['light', 'dark', 'system']
    const currentIndex = themes.indexOf(theme.value)
    const nextIndex = (currentIndex + 1) % themes.length
    setTheme(themes[nextIndex])
  }

  // Initialize on client
  onMounted(() => {
    // Load from localStorage
    const stored = localStorage.getItem('theme') as Theme | null
    if (stored && ['light', 'dark', 'system'].includes(stored)) {
      theme.value = stored
    }
    updateResolvedTheme()

    // Listen for system theme changes
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
    mediaQuery.addEventListener('change', updateResolvedTheme)

    onUnmounted(() => {
      mediaQuery.removeEventListener('change', updateResolvedTheme)
    })
  })

  return {
    theme: resolvedTheme,
    rawTheme: theme,
    setTheme,
    toggleTheme,
  }
}

