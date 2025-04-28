import { jwtDecode } from 'jwt-decode'

interface DecodeToken {
  exp: number
}

export function saveTokens(accessToken: string, refreshToken: string): void {
  sessionStorage.setItem('access_token', accessToken)
  sessionStorage.setItem('refresh_token', refreshToken)
}

export function getAccessToken(): string | null {
  if (typeof window !== 'undefined' && typeof sessionStorage !== 'undefined') {
    return sessionStorage.getItem('access_token')
  }
  return null
}

export function getRefreshToken(): string | null {
  if (typeof window !== 'undefined' && typeof sessionStorage !== 'undefined') {
    return sessionStorage.getItem('refresh_token')
  }
  return null
}

export function isTokenExpired(token: string): boolean {
  const decoded: DecodeToken = jwtDecode(token)
  const curTime = Math.floor(Date.now() / 1000)

  return decoded.exp < curTime
}
