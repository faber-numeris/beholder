import { createBrowserRouter, Navigate, Outlet } from 'react-router'
import { useAuth } from './auth'
import Login from './pages/Login'
import Signup from './pages/Signup'
import Dashboard from './pages/Dashboard'

function ProtectedRoute() {
  const { authenticated } = useAuth()
  if (!authenticated) return <Navigate to="/" replace />
  return <Outlet />
}

function PublicRoute() {
  const { authenticated } = useAuth()
  if (authenticated) return <Navigate to="/dashboard" replace />
  return <Outlet />
}

export const router = createBrowserRouter([
  {
    element: <PublicRoute />,
    children: [
      { path: '/', element: <Login /> },
      { path: '/signup', element: <Signup /> },
    ],
  },
  {
    element: <ProtectedRoute />,
    children: [
      { path: '/dashboard', element: <Dashboard /> },
    ],
  },
])
