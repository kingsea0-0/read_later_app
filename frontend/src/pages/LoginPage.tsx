import { useAuth } from '../contexts/AuthContext'
import { useNavigate } from 'react-router-dom'
import { useEffect } from 'react'

export default function LoginPage() {
  const { user, loading, signIn } = useAuth()
  const navigate = useNavigate()

  useEffect(() => {
    if (user) navigate('/', { replace: true })
  }, [user, navigate])

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-neutral-50">
        <p className="text-neutral-500">Loading...</p>
      </div>
    )
  }

  return (
    <div className="flex items-center justify-center min-h-screen bg-neutral-50">
      <div className="w-full max-w-sm mx-4">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-neutral-900">Hoard</h1>
          <p className="mt-2 text-neutral-500">Save links. Read later.</p>
        </div>
        <div className="space-y-3">
          <button
            onClick={() => signIn('google')}
            className="w-full px-4 py-3 text-sm font-medium text-neutral-700 bg-white border border-neutral-300 rounded-lg hover:bg-neutral-50 transition-colors"
          >
            Continue with Google
          </button>
          <button
            onClick={() => signIn('apple')}
            className="w-full px-4 py-3 text-sm font-medium text-white bg-neutral-900 rounded-lg hover:bg-neutral-800 transition-colors"
          >
            Continue with Apple
          </button>
        </div>
      </div>
    </div>
  )
}
