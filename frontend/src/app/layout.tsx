import type { Metadata } from 'next'
import { Geist, Geist_Mono } from 'next/font/google'
import './globals.css'
import BackgroundWrapper from '@/components/BackgroundWrapper'

const geistSans = Geist({
  variable: '--font-geist-sans',
  subsets: ['latin']
})

const geistMono = Geist_Mono({
  variable: '--font-geist-mono',
  subsets: ['latin']
})

export const metadata: Metadata = {
  title: 'LDAP service',
  description: 'LDAP service DNP project UI',
  authors: [
    { name: 'Projacktor', url: 'https://github.com/projacktor' },
    { name: 'Woolfer0097', url: 'https://github.com/Woolfer0097' },
    { name: 'Abraham14711', url: 'https://github.com/Abraham14711' },
    { name: 'Pickpusha', url: 'https://github.com/ilyalinhnguyen' },
    { name: 'FunnyFoXD', url: 'https://github.com/FunnyFoXD' }
  ]
}

export default function RootLayout({
  children
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="en">
      <body className={`${geistSans.variable} ${geistMono.variable} bg-black antialiased`}>
        <BackgroundWrapper />
        {children}
      </body>
    </html>
  )
}
