'use client'
import { User } from '@/types/types'
import Header from '@/components/Header'
import { Button } from '@/components/ui/button'
import Clocks from '@/components/clocks'
import { logout } from '@/lib/auth'
import { getUserInfo } from '@/lib/api'
import { useEffect, useState } from 'react'

export default function Home() {
  const [user, setUser] = useState<User | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    async function fetchUser() {
      try {
        const userInfo = await getUserInfo()
        setUser(userInfo)
      } catch (err) {
        console.error(err)
        setError('Failed to connect to server')
      }
    }

    fetchUser()
  }, [])

  return (
    <>
      <Header />
      <main className="h-full">
        {user ? (
          <div className="flex flex-col items-center space-y-8">
            <Clocks />

            <section className="space-y-4 rounded-4xl border-2 border-gray-400 bg-white/3 p-10 text-left text-white backdrop-blur-sm">
              <div className="flex flex-row items-center space-x-3">
                <h2 className="text-3xl font-semibold">Username:</h2>
                <h3 className="text-2xl font-normal hover:underline">{user.preferred_username}</h3>
              </div>
              <div className="flex flex-row items-center space-x-3">
                <h2 className="text-3xl font-semibold">Your email:</h2>
                <h3 className="text-2xl font-normal hover:underline">{user.email}</h3>
              </div>
              <div className="flex flex-row items-center space-x-3">
                <h2 className="text-3xl font-semibold">Your roles:</h2>
                <h3 className="text-2xl font-normal hover:underline">
                  {user.resource_access.account.roles.join(', ')}
                </h3>
              </div>
            </section>
            <Button variant={'secondary'} className="h-10 w-90" onClick={logout}>
              Log out
            </Button>
          </div>
        ) : (
          <p className="text-center text-xl font-semibold text-[#c41010]">
            {error || 'Loading...'}
          </p>
        )}
      </main>
    </>
  )
}
