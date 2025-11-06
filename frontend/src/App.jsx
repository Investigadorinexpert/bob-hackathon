import { useState } from 'react'
import ChatWidget from './components/ChatWidget'
import LeadsDashboard from './components/LeadsDashboard'
import './App.css'

function App() {
  const [activeView, setActiveView] = useState('chat')

  return (
    <div className="app">
      <header className="app-header">
        <div className="header-content">
          <h1>🤖 BOB Chatbot</h1>
          <p>Asistente Virtual de Subastas</p>
        </div>
        <nav className="nav-tabs">
          <button
            className={`nav-tab ${activeView === 'chat' ? 'active' : ''}`}
            onClick={() => setActiveView('chat')}
          >
            💬 Chat
          </button>
          <button
            className={`nav-tab ${activeView === 'leads' ? 'active' : ''}`}
            onClick={() => setActiveView('leads')}
          >
            📊 Leads
          </button>
        </nav>
      </header>

      <main className="app-main">
        {activeView === 'chat' ? <ChatWidget /> : <LeadsDashboard />}
      </main>

      <footer className="app-footer">
        <p>Hackathon BOB 2025 | Powered by Gemini AI 2.5 Flash</p>
      </footer>
    </div>
  )
}

export default App
