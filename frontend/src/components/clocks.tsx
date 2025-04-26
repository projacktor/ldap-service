<<<<<<< HEAD
import React from 'react'

function Clock() {
  return (
    <div>
      <Clock
        format={'h:mm'ss}
      />
=======
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
>>>>>>> f990bed1d29dbf7000550e38b459f28fec178b89
    </div>
  )
}

<<<<<<< HEAD
export default Clock
=======
export default Clocks
>>>>>>> f990bed1d29dbf7000550e38b459f28fec178b89
