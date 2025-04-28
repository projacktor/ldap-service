import { jwtDecode } from 'jwt-decode'
import Cookies from 'js-cookie'

interface DecodeToken {
  exp: number
}

export function saveTokens(accessToken: string, refreshToken: string): void {
  Cookies.set('access_token', accessToken, {
    expires: 1 / 24,
    secure: true,
    sameSite: 'Strict'
  })
  Cookies.set('refresh_token', refreshToken, {
    expires: 7,
    secure: true,
    sameSite: 'Strict'
  })
}

export function getAccessToken(): string | null {
  return Cookies.get('access_token') || null
}

export function getRefreshToken(): string | null {
  return Cookies.get('refresh_token') || null
}
export function isTokenExpired(token: string): boolean {
  const decoded: DecodeToken = jwtDecode(token)
  const curTime = Math.floor(Date.now() / 1000)

  return decoded.exp < curTime
}

export function logout(): void {
  Cookies.remove('access_token')
  Cookies.remove('refresh_token')
  window.location.href = '/auth'
}
