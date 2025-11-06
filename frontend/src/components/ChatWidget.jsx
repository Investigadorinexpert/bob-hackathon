import { useState, useRef, useEffect } from 'react'
import './ChatWidget.css'

function ChatWidget() {
  const [messages, setMessages] = useState([
    {
      role: 'assistant',
      content: '¡Hola! Soy tu asistente de BOB Subastas. ¿En qué puedo ayudarte hoy?',
      timestamp: new Date().toISOString()
    }
  ])
  const [inputMessage, setInputMessage] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [sessionId, setSessionId] = useState(null)
  const [leadScore, setLeadScore] = useState(0)
  const [category, setCategory] = useState('cold')
  const messagesEndRef = useRef(null)

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }

  useEffect(() => {
    scrollToBottom()
  }, [messages])

  const sendMessage = async () => {
    if (!inputMessage.trim() || isLoading) return

    const userMessage = inputMessage.trim()
    setInputMessage('')
    setIsLoading(true)

    // Agregar mensaje del usuario
    const newUserMessage = {
      role: 'user',
      content: userMessage,
      timestamp: new Date().toISOString()
    }
    setMessages(prev => [...prev, newUserMessage])

    try {
      const response = await fetch('/api/chat/message', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          sessionId,
          message: userMessage,
          channel: 'web'
        })
      })

      const data = await response.json()

      if (data.success) {
        // Guardar sessionId si es nuevo
        if (!sessionId) {
          setSessionId(data.sessionId)
        }

        // Actualizar score
        setLeadScore(data.leadScore || 0)
        setCategory(data.category || 'cold')

        // Agregar respuesta del asistente
        const assistantMessage = {
          role: 'assistant',
          content: data.reply,
          timestamp: data.timestamp
        }
        setMessages(prev => [...prev, assistantMessage])
      } else {
        throw new Error(data.error || 'Error desconocido')
      }
    } catch (error) {
      console.error('Error sending message:', error)
      const errorMessage = {
        role: 'assistant',
        content: 'Lo siento, hubo un error al procesar tu mensaje. Por favor intenta de nuevo.',
        timestamp: new Date().toISOString(),
        isError: true
      }
      setMessages(prev => [...prev, errorMessage])
    } finally {
      setIsLoading(false)
    }
  }

  const handleKeyPress = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      sendMessage()
    }
  }

  const getCategoryColor = () => {
    switch (category) {
      case 'hot': return '#f56565'
      case 'warm': return '#ed8936'
      case 'cold': return '#4299e1'
      default: return '#cbd5e0'
    }
  }

  const getCategoryEmoji = () => {
    switch (category) {
      case 'hot': return '🔥'
      case 'warm': return '😐'
      case 'cold': return '❄️'
      default: return '❓'
    }
  }

  return (
    <div className="chat-widget">
      <div className="chat-score-panel">
        <div className="score-card">
          <div className="score-label">Lead Score</div>
          <div className="score-value" style={{ color: getCategoryColor() }}>
            {leadScore}
          </div>
          <div className="score-category">
            {getCategoryEmoji()} {category.toUpperCase()}
          </div>
        </div>
        {sessionId && (
          <div className="session-info">
            <small>Session: {sessionId.slice(0, 20)}...</small>
          </div>
        )}
      </div>

      <div className="chat-container">
        <div className="chat-messages">
          {messages.map((msg, index) => (
            <div
              key={index}
              className={`message ${msg.role} ${msg.isError ? 'error' : ''}`}
            >
              <div className="message-content">
                {msg.content}
              </div>
              <div className="message-time">
                {new Date(msg.timestamp).toLocaleTimeString('es-PE', {
                  hour: '2-digit',
                  minute: '2-digit'
                })}
              </div>
            </div>
          ))}
          {isLoading && (
            <div className="message assistant">
              <div className="message-content">
                <div className="typing-indicator">
                  <span></span>
                  <span></span>
                  <span></span>
                </div>
              </div>
            </div>
          )}
          <div ref={messagesEndRef} />
        </div>

        <div className="chat-input-container">
          <textarea
            value={inputMessage}
            onChange={(e) => setInputMessage(e.target.value)}
            onKeyPress={handleKeyPress}
            placeholder="Escribe tu mensaje..."
            rows="2"
            disabled={isLoading}
          />
          <button
            onClick={sendMessage}
            disabled={!inputMessage.trim() || isLoading}
            className="send-button"
          >
            {isLoading ? '⏳' : '📤'}
          </button>
        </div>
      </div>
    </div>
  )
}

export default ChatWidget
