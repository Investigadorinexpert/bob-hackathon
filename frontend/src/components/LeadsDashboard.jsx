import { useState, useEffect } from 'react'
import './LeadsDashboard.css'

function LeadsDashboard() {
  const [leads, setLeads] = useState([])
  const [stats, setStats] = useState(null)
  const [isLoading, setIsLoading] = useState(true)
  const [filter, setFilter] = useState('all')

  useEffect(() => {
    fetchLeads()
    fetchStats()
    const interval = setInterval(() => {
      fetchLeads()
      fetchStats()
    }, 5000) // Actualizar cada 5 segundos

    return () => clearInterval(interval)
  }, [filter])

  const fetchLeads = async () => {
    try {
      const queryParams = filter !== 'all' ? `?category=${filter}` : ''
      const response = await fetch(`/api/leads${queryParams}`)
      const data = await response.json()

      if (data.success) {
        setLeads(data.leads)
      }
      setIsLoading(false)
    } catch (error) {
      console.error('Error fetching leads:', error)
      setIsLoading(false)
    }
  }

  const fetchStats = async () => {
    try {
      const response = await fetch('/api/leads/stats')
      const data = await response.json()

      if (data.success) {
        setStats(data.stats)
      }
    } catch (error) {
      console.error('Error fetching stats:', error)
    }
  }

  const getCategoryColor = (category) => {
    switch (category) {
      case 'hot': return '#f56565'
      case 'warm': return '#ed8936'
      case 'cold': return '#4299e1'
      default: return '#cbd5e0'
    }
  }

  const getCategoryEmoji = (category) => {
    switch (category) {
      case 'hot': return '🔥'
      case 'warm': return '😐'
      case 'cold': return '❄️'
      default: return '❓'
    }
  }

  const getChannelEmoji = (channel) => {
    switch (channel) {
      case 'web': return '💻'
      case 'whatsapp': return '💬'
      default: return '📱'
    }
  }

  if (isLoading) {
    return <div className="loading">Cargando leads...</div>
  }

  return (
    <div className="leads-dashboard">
      {stats && (
        <div className="stats-grid">
          <div className="stat-card">
            <div className="stat-label">Total Leads</div>
            <div className="stat-value">{stats.total}</div>
          </div>
          <div className="stat-card hot">
            <div className="stat-label">🔥 Calientes</div>
            <div className="stat-value">{stats.hot}</div>
          </div>
          <div className="stat-card warm">
            <div className="stat-label">😐 Tibios</div>
            <div className="stat-value">{stats.warm}</div>
          </div>
          <div className="stat-card cold">
            <div className="stat-label">❄️ Fríos</div>
            <div className="stat-value">{stats.cold}</div>
          </div>
          <div className="stat-card">
            <div className="stat-label">Score Promedio</div>
            <div className="stat-value">{stats.avgScore.toFixed(0)}</div>
          </div>
        </div>
      )}

      <div className="filter-bar">
        <button
          className={`filter-btn ${filter === 'all' ? 'active' : ''}`}
          onClick={() => setFilter('all')}
        >
          Todos
        </button>
        <button
          className={`filter-btn ${filter === 'hot' ? 'active' : ''}`}
          onClick={() => setFilter('hot')}
        >
          🔥 Calientes
        </button>
        <button
          className={`filter-btn ${filter === 'warm' ? 'active' : ''}`}
          onClick={() => setFilter('warm')}
        >
          😐 Tibios
        </button>
        <button
          className={`filter-btn ${filter === 'cold' ? 'active' : ''}`}
          onClick={() => setFilter('cold')}
        >
          ❄️ Fríos
        </button>
      </div>

      <div className="leads-list">
        {leads.length === 0 ? (
          <div className="no-leads">
            <p>No hay leads aún. ¡Empieza una conversación en el chat!</p>
          </div>
        ) : (
          leads.map((lead) => (
            <div key={lead.sessionId} className="lead-card">
              <div className="lead-header">
                <div className="lead-id">
                  {getChannelEmoji(lead.channel)} {lead.sessionId.slice(0, 25)}...
                </div>
                <div
                  className="lead-score"
                  style={{ color: getCategoryColor(lead.category) }}
                >
                  {getCategoryEmoji(lead.category)} {lead.score}
                </div>
              </div>

              <div className="lead-details">
                <div className="detail-item">
                  <span className="detail-label">Categoría:</span>
                  <span className="detail-value" style={{ color: getCategoryColor(lead.category) }}>
                    {lead.category.toUpperCase()}
                  </span>
                </div>
                <div className="detail-item">
                  <span className="detail-label">Urgencia:</span>
                  <span className="detail-value">{lead.urgency || 'unknown'}</span>
                </div>
                <div className="detail-item">
                  <span className="detail-label">Presupuesto:</span>
                  <span className="detail-value">{lead.budget || 'unknown'}</span>
                </div>
                <div className="detail-item">
                  <span className="detail-label">Tipo:</span>
                  <span className="detail-value">{lead.businessType || 'unknown'}</span>
                </div>
                <div className="detail-item">
                  <span className="detail-label">Mensajes:</span>
                  <span className="detail-value">{lead.messageCount}</span>
                </div>
              </div>

              {lead.lastMessage && (
                <div className="last-message">
                  <strong>Último mensaje:</strong> {lead.lastMessage}
                </div>
              )}

              <div className="lead-footer">
                <small>
                  Actualizado: {new Date(lead.updatedAt).toLocaleString('es-PE')}
                </small>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  )
}

export default LeadsDashboard
