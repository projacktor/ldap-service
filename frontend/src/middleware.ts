import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'
import { isTokenExpired } from '@/lib/auth'
import { refreshToken } from '@/lib/api'

export async function middleware(req: NextRequest) {
  const accessToken = req.cookies.get('access_token')?.value

  if (!accessToken || isTokenExpired(accessToken)) {
    try {
      const newAccessToken = await refreshToken()
      const res = NextResponse.next()
      res.cookies.set('access_token', newAccessToken, { httpOnly: true })
      return res
    } catch (err) {
      return NextResponse.redirect(new URL('/auth', req.url))
    }
  }

  return NextResponse.next()
}

export const config = {
  matcher: ['/']
}
