import React from 'react'
import { Input } from '@/components/ui/input'

function Page() {
  return (
    <main className="content-center">
      <form className="flex flex-col content-center space-y-4">
        <h1 className="text-center text-4xl font-bold text-white">Sign in to your account</h1>
        <p className="text-center text-base font-normal text-gray-500">
          Enter your email below to enter to your account
        </p>
        <Input type="username" placeholder="Enter your username" />

      </form>
    </main>
  )
}

export default Page
