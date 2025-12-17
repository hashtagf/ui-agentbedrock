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

interface Message {
  id: string
  sessionId: string
  role: 'user' | 'assistant'
  content: string
  trace?: any
  createdAt: string
}

export function useChat() {
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase
  
  const { currentSession, messages, addMessage, updateLastMessage, clearMessages } = useSession()
  
  const isStreaming = useState('isStreaming', () => false)
  const thinkingStatus = useState<string | null>('thinkingStatus', () => null)
  const agentSteps = useState<AgentStep[]>('agentSteps', () => [])
  const currentError = useState<any>('currentError', () => null)
  const abortController = useState<AbortController | null>('abortController', () => null)
  const wasSummarized = useState<boolean>('wasSummarized', () => false)

  const sendMessage = async (content: string) => {
    if (!currentSession.value || isStreaming.value || !content.trim()) return

    // Reset state
    thinkingStatus.value = null
    agentSteps.value = []
    currentError.value = null
    wasSummarized.value = false

    // Add user message
    const userMessage: Message = {
      id: `temp-${Date.now()}`,
      sessionId: currentSession.value.id,
      role: 'user',
      content: content.trim(),
      createdAt: new Date().toISOString(),
    }
    addMessage(userMessage)

    // Add placeholder assistant message
    const assistantMessage: Message = {
      id: `temp-assistant-${Date.now()}`,
      sessionId: currentSession.value.id,
      role: 'assistant',
      content: '',
      createdAt: new Date().toISOString(),
    }
    addMessage(assistantMessage)

    // Start streaming
    isStreaming.value = true
    abortController.value = new AbortController()

    try {
      const response = await fetch(`${apiBase}/api/chat/stream`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          sessionId: currentSession.value.id,
          message: content.trim(),
        }),
        signal: abortController.value.signal,
      })

      if (!response.ok) {
        throw new Error('Failed to send message')
      }

      const reader = response.body?.getReader()
      const decoder = new TextDecoder()
      let buffer = ''
      let fullContent = ''

      while (reader) {
        const { done, value } = await reader.read()
        if (done) break

        buffer += decoder.decode(value, { stream: true })
        const lines = buffer.split('\n')
        buffer = lines.pop() || ''

        for (const line of lines) {
          if (line.startsWith('event:')) {
            const eventType = line.slice(6).trim()
            continue
          }
          
          if (line.startsWith('data:')) {
            const dataStr = line.slice(5).trim()
            if (!dataStr) continue

            try {
              const data = JSON.parse(dataStr)
              
              // Determine event type from previous line or data structure
              if (data.status !== undefined && typeof data.status === 'string' && !data.stepIndex) {
                // thinking event
                thinkingStatus.value = data.status
              } else if (data.stepIndex !== undefined) {
                // agent_step event
                const stepIndex = agentSteps.value.findIndex(s => s.stepIndex === data.stepIndex)
                if (stepIndex === -1) {
                  agentSteps.value = [...agentSteps.value, data as AgentStep]
                } else {
                  agentSteps.value[stepIndex] = data as AgentStep
                }
              } else if (data.chunk !== undefined) {
                // content event
                fullContent += data.chunk
                updateLastMessage(fullContent)
                thinkingStatus.value = null
              } else if (data.traceId !== undefined) {
                // trace event
                updateLastMessage(fullContent, data)
              } else if (data.type !== undefined && data.message !== undefined && data.type !== 'StreamError') {
                // error event
                currentError.value = data
              } else if (data.message !== undefined && !data.type) {
                // summarized event
                wasSummarized.value = true
              } else if (data.messageId !== undefined) {
                // done event
                thinkingStatus.value = null
              }
            } catch (e) {
              console.error('Failed to parse SSE data:', e)
            }
          }
        }
      }
    } catch (error: any) {
      if (error.name !== 'AbortError') {
        console.error('Stream error:', error)
        currentError.value = {
          type: 'StreamError',
          message: error.message,
        }
      }
    } finally {
      isStreaming.value = false
      abortController.value = null
    }
  }

  const stopStream = () => {
    if (abortController.value) {
      abortController.value.abort()
      isStreaming.value = false
      thinkingStatus.value = null
    }
  }

  const clearError = () => {
    currentError.value = null
  }

  const clearHistory = async () => {
    if (!currentSession.value) return
    
    try {
      const response = await fetch(`${apiBase}/api/sessions/${currentSession.value.id}/messages`, {
        method: 'DELETE',
      })
      
      if (!response.ok) {
        throw new Error('Failed to clear history')
      }
      
      // Clear local messages
      clearMessages()
      wasSummarized.value = false
    } catch (error: any) {
      console.error('Failed to clear history:', error)
      currentError.value = {
        type: 'ClearHistoryError',
        message: error.message,
      }
    }
  }

  return {
    messages,
    isStreaming,
    thinkingStatus,
    agentSteps,
    currentError,
    wasSummarized,
    sendMessage,
    stopStream,
    clearError,
    clearHistory,
  }
}

