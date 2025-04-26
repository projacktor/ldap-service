'use client'
import React, { useEffect, useState } from 'react'
import Clock from 'react-live-clock'

function Clocks() {
  const [isMounted, setIsMounted] = useState(false)

  useEffect(() => {
    setIsMounted(true)
  }, [])

  if (!isMounted) return null

  return (
    <div className="rounded-4xl border-2 border-gray-400 p-10 text-white">
      <Clock format={'h:mm:ss a'} style={{ fontSize: '2rem' }} ticking={true} />
    </div>
  )
}

export default Clocks
