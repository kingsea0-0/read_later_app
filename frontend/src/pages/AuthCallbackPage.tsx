import { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { supabase } from '../lib/supabase'

export default function AuthCallbackPage() {
  const navigate = useNavigate()

  useEffect(() => {
    supabase.auth.getSession().then(({ data: { session } }) => {
      if (session) {
        // Notify Chrome Extension via postMessage
        window.postMessage({
          type: 'HOARD_LOGIN',
          accessToken: session.access_token,
          refreshToken: session.refresh_token,
        }, window.location.origin)

        navigate('/', { replace: true })
      } else {
        navigate('/auth', { replace: true })
      }
    })
  }, [navigate])

  return (
    <div className="flex items-center justify-center min-h-screen bg-neutral-50">
      <p className="text-neutral-500">Completing sign in...</p>
    </div>
  )
}
