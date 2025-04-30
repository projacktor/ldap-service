import type { NextConfig } from 'next'
import * as dotenv from 'dotenv'

const nextConfig: NextConfig = {
  /* config options here */
  env: {
    NEXT_PUBLIC_BACKEND_IP: process.env.NEXT_PUBLIC_BACKEND_IP
  }
}

export default nextConfig
