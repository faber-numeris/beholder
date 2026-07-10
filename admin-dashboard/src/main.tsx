import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'

import 'primereact/resources/themes/lara-light-amber/theme.css'
import 'primeicons/primeicons.css'
import 'primeflex/primeflex.css'
import './index.css'

import { AuthProvider } from './auth'
import App from './App.tsx'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <AuthProvider>
      <App />
    </AuthProvider>
  </StrictMode>,
)
