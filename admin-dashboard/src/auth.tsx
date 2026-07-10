import { createContext, useContext, useState, useCallback, type ReactNode } from 'react'

export interface AuthContext {
  authenticated: boolean
  login: () => void
  logout: () => void
}

const AuthCtx = createContext<AuthContext | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [authenticated, setAuthenticated] = useState(false)

  const login = useCallback(() => setAuthenticated(true), [])
  const logout = useCallback(() => setAuthenticated(false), [])

  return (
    <AuthCtx.Provider value={{ authenticated, login, logout }}>
      {children}
    </AuthCtx.Provider>
  )
}

export function useAuth(): AuthContext {
  const ctx = useContext(AuthCtx)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
