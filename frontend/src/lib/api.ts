import { getRefreshToken, saveTokens } from '@/lib/auth'

export async function verifyUser(username: string, password: string) {
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_IP}/auth/login`, {
    method: 'POST',
    headers: {
      'Content-type': 'application/json'
    },
    body: JSON.stringify({
      username: username,
      password: password
    })
  })

  if (!response.ok) {
    throw new Error('Failed to validate user')
  }

  const data = await response.json()
  saveTokens(data.access_token, data.refresh_token)
  return data
}

export async function refreshToken(): Promise<string> {
  const refreshToken = getRefreshToken()
  if (!refreshToken) {
    throw new Error('Refresh token not found')
  }

  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_IP}/auth/refresh`, {
    method: 'POST',
    headers: {
      'Content-type': 'application/json'
    },
    body: JSON.stringify({ refreshToken: refreshToken })
  })

  if (!response.ok) {
    throw new Error('Failed to refresh token')
  }

  const data = await response.json()
  saveTokens(data.access_token, data.refreshToken)
  return data.access_token
}
