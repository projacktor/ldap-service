import { jwtDecode } from 'jwt-decode'
import Cookies from 'js-cookie'

interface DecodeToken {
  exp: number
}

export function saveTokens(
  accessToken: string,
  refreshToken: string,
  idToken: string,
  expiresIn: number,
  refreshExpiresIn: number
): void {
  const accessTokenExp = new Date(Date.now() + expiresIn * 1000)
  const refreshTokenExp = new Date(Date.now() + refreshExpiresIn * 1000)

  Cookies.set('access_token', accessToken, {
    expires: accessTokenExp,
    secure: true,
    sameSite: 'Strict'
  })
  Cookies.set('refresh_token', refreshToken, {
    expires: refreshTokenExp,
    secure: true,
    sameSite: 'Strict'
  })
  Cookies.set('id_token', idToken, {
    expires: 10e9, // long value
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

export function getIdToken(): string | null {
  return Cookies.get('id_token') || null
}

export function isTokenExpired(token: string): boolean {
  const decoded: DecodeToken = jwtDecode(token)
  const curTime = Math.floor(Date.now() / 1000)

  return decoded.exp < curTime
}

export function logout(): void {
  Cookies.remove('access_token')
  Cookies.remove('refresh_token')
  Cookies.remove('id_token')
  window.location.href = '/auth'
}
