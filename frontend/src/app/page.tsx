'use client'
import { User } from '@/lib/types'
import Header from '@/components/Header'
import { useRef } from 'react'
import { Button } from '@/components/ui/button'
import Clocks from '@/components/clocks'

export default function Home() {
  const user: User = {
    username: 'USSSER',
    role: 'SUPERPUPER'
  }

  const mainContRef = useRef<HTMLDivElement>(null)
  return (
    <>
      <Header />
      <main className="h-full">
        {user ? (
          <div className="flex flex-col items-center space-y-8">
            <Clocks />

            <section className="space-y-4 rounded-4xl border-2 border-gray-400 p-10 text-left text-white">
              <div className="flex flex-row items-center space-x-3">
                <h2 className="text-3xl font-semibold">Username:</h2>
                <h3 className="text-2xl font-normal hover:underline">{user.username}</h3>
              </div>
              <div className="flex flex-row items-center space-x-3">
                <h2 className="text-3xl font-semibold">Your role:</h2>
                <h3 className="text-2xl font-normal hover:underline">{user.role}</h3>
              </div>
            </section>
            <Button variant={'secondary'} className="h-10 w-90">
              Log out
            </Button>
          </div>
        ) : (
          <p className="text-center text-lg font-semibold text-[#c41010]">
            Failed to connect to server
          </p>
        )}
      </main>
    </>
  )
}
