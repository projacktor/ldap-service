import { getIdToken, getRefreshToken, isTokenExpired, saveTokens } from '@/lib/auth'
import { User } from '@/lib/types'

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
  saveTokens(
    data.access_token,
    data.refresh_token,
    data.id_token,
    data.expires_in,
    data.refresh_expires_in
  )
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
  saveTokens(
    data.access_token,
    data.refreshToken,
    data.id_token,
    data.expires_in,
    data.refresh_expires_in
  )
  return data.access_token
}

export async function getUserInfo(): Promise<User> {
  const refreshToken = getRefreshToken()
  if (refreshToken !== null && !isTokenExpired(refreshToken)) {
    const idToken = getIdToken()

    if (!idToken) {
      throw new Error('Id token not found')
    }

    const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_IP}/users`, {
      method: 'GET',
      headers: {
        'Content-type': 'application/json',
        Authorization: `Bearer ${idToken}`
      }
    })

    if (!response.ok) {
      throw new Error('Failed to get user info')
    }
    return (await response.json()) as User
  } else {
    throw new Error('Refresh token expired or not found')
  }
}
