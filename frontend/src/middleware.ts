import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'
import { getAccessToken, isTokenExpired } from '@/lib/auth'
import { refreshToken } from '@/lib/api'

export async function middleware(req: NextRequest) {
  const accessToken = getAccessToken()

  if (!accessToken || isTokenExpired(accessToken)) {
    try {
      const newAccessToken = await refreshToken()
      req.headers.set('Authorization', `Bearer ${newAccessToken}`)
    } catch (err) {
      return NextResponse.redirect(new URL('/auth', req.url))
    }
  }

  return NextResponse.next()
}

export const config = {
  matcher: ['/']
}
