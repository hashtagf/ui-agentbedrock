interface Session {
  id: string
  title: string
  createdAt: string
  updatedAt: string
}

interface Message {
  id: string
  sessionId: string
  role: 'user' | 'assistant'
  content: string
  trace?: Trace
  createdAt: string
}

interface Trace {
  traceId: string
  agentSteps: AgentStep[]
  error?: ErrorInfo
}

interface AgentStep {
  stepIndex: number
  agentName: string
  action: string
  status: 'running' | 'success' | 'error'
  startTime: string
  endTime?: string
}

interface ErrorInfo {
  type: string
  message: string
  source: string
  stackTrace?: string
}

export function useSession() {
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase

  const sessions = useState<Session[]>('sessions', () => [])
  const currentSession = useState<Session | null>('currentSession', () => null)
  const currentMessages = useState<Message[]>('currentMessages', () => [])
  const isLoading = useState('sessionsLoading', () => false)

  const fetchSessions = async () => {
    isLoading.value = true
    try {
      const response = await fetch(`${apiBase}/api/sessions`)
      if (response.ok) {
        sessions.value = await response.json()
      }
    } catch (error) {
      console.error('Failed to fetch sessions:', error)
    } finally {
      isLoading.value = false
    }
  }

  const createSession = async (title?: string) => {
    try {
      const response = await fetch(`${apiBase}/api/sessions`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title: title || 'New Chat' }),
      })
      
      if (response.ok) {
        const session = await response.json()
        sessions.value = [session, ...sessions.value]
        await selectSession(session.id)
      }
    } catch (error) {
      console.error('Failed to create session:', error)
    }
  }

  const selectSession = async (sessionId: string) => {
    try {
      const response = await fetch(`${apiBase}/api/sessions/${sessionId}`)
      if (response.ok) {
        const data = await response.json()
        currentSession.value = data.session
        currentMessages.value = data.messages || []
      }
    } catch (error) {
      console.error('Failed to fetch session:', error)
    }
  }

  const updateSession = async (sessionId: string, title: string) => {
    try {
      const response = await fetch(`${apiBase}/api/sessions/${sessionId}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title }),
      })
      
      if (response.ok) {
        const index = sessions.value.findIndex(s => s.id === sessionId)
        if (index !== -1) {
          sessions.value[index].title = title
        }
        if (currentSession.value?.id === sessionId) {
          currentSession.value.title = title
        }
      }
    } catch (error) {
      console.error('Failed to update session:', error)
    }
  }

  const deleteSession = async (sessionId: string) => {
    try {
      const response = await fetch(`${apiBase}/api/sessions/${sessionId}`, {
        method: 'DELETE',
      })
      
      if (response.ok) {
        sessions.value = sessions.value.filter(s => s.id !== sessionId)
        if (currentSession.value?.id === sessionId) {
          currentSession.value = null
          currentMessages.value = []
        }
      }
    } catch (error) {
      console.error('Failed to delete session:', error)
    }
  }

  const addMessage = (message: Message) => {
    currentMessages.value = [...currentMessages.value, message]
  }

  const updateLastMessage = (content: string, trace?: Trace) => {
    if (currentMessages.value.length > 0) {
      const lastIndex = currentMessages.value.length - 1
      const lastMessage = currentMessages.value[lastIndex]
      currentMessages.value[lastIndex] = {
        ...lastMessage,
        content,
        trace: trace || lastMessage.trace,
      }
    }
  }

  const clearMessages = () => {
    currentMessages.value = []
  }

  // Initialize
  onMounted(() => {
    fetchSessions()
  })

  return {
    sessions,
    currentSession,
    messages: currentMessages,
    isLoading,
    fetchSessions,
    createSession,
    selectSession,
    updateSession,
    deleteSession,
    addMessage,
    updateLastMessage,
    clearMessages,
  }
}

