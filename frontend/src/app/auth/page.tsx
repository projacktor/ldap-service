'use client'
import React from 'react'
import { Input } from '@/components/ui/input'
import { useForm, SubmitHandler } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Button } from '@/components/ui/button'

const formSchema = z.object({
  username: z.string().min(4, {
    message: 'Username must be at least 4 characters.'
  }),
  password: z.string().min(8, {
    message: 'Password must be at least 8 characters.'
  })
})

type FormSchema = z.infer<typeof formSchema>

function Page() {
  const {
    register,
    handleSubmit,
    formState: { errors }
  } = useForm<FormSchema>({
    resolver: zodResolver(formSchema)
  })

  const onSubmit: SubmitHandler<FormSchema> = (data) => console.log(data)

  return (
    <main className="h-screen">
      <div className="flex flex-col items-center space-y-8">
        <article className="space-y-4">
          <h1 className="text-center text-5xl font-bold text-white">Sign in to your account</h1>
          <p className="text-center text-base font-normal text-gray-400">
            Enter your email below to enter to your account
          </p>
        </article>
        <form
          onSubmit={handleSubmit(onSubmit)}
          className="flex flex-col items-center space-y-5 text-white"
        >
          <div className="w-90 space-y-2">
            {errors.username && (
              <p className="text-left text-sm text-[#c41010]">{errors.username.message}</p>
            )}
            <Input
              id="username"
              type="text"
              {...register('username')}
              placeholder="Username"
              className="h-10 text-white"
            />
          </div>

          <div className="w-90 space-y-2">
            {errors.password && <p className="text-sm text-[#c41010]">{errors.password.message}</p>}
            <Input
              id="password"
              type="password"
              {...register('password')}
              placeholder="Password"
              className="h-10"
            />
          </div>
          <Button type="submit" variant="secondary" className="h-10 w-90">
            Sign in
          </Button>
        </form>
      </div>
    </main>
  )
}

export default Page
