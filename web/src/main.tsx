import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import 'normalize.css/normalize.css'
import './index.css'
import App from './App.tsx'
import { BrowserRouter } from 'react-router-dom'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <BrowserRouter basename="/web">
      <App />
    </BrowserRouter>
  </StrictMode>,
)
